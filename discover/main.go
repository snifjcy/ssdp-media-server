package main

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"golang.org/x/sys/windows"
)

// func main() {
// 	fmt.Printf("%+v", Discover(3, "ssdp:all"))
// }

func main() {
	mod, err := windows.LoadDLL("AgentModule.dll")
	if err != nil {
		panic(err)
	}
	_, err = mod.FindProc("MakeAuthorizationHeaderWithGeneratedNonceValueAndAMModule")
	if err != nil {
		panic(err)
	}
}

// buildSearch returns the message used to discover devices
// mx: Maximum wait time (in seconds)
// st: Search target
func buildSearch(mx int, st string) []byte {
	discoverMsg := &bytes.Buffer{}
	fmt.Fprint(discoverMsg, "M-SEARCH * HTTP/1.1\r\n")
	fmt.Fprint(discoverMsg, "HOST: 239.255.255.250:1900\r\n")
	fmt.Fprint(discoverMsg, "MAN: \"ssdp:discover\"\r\n")
	fmt.Fprintf(discoverMsg, "MX: %d\r\n", mx)
	fmt.Fprintf(discoverMsg, "ST: %s\r\n\r\n", st)
	return discoverMsg.Bytes()
}

// Discover function
func Discover(timeout int, target string) []string {
	var searchResults []string

	if timeout <= 2 {
		timeout = 2
	}

	go func() {
		mx := timeout - 1

		conn, err := net.ListenUDP("udp", nil)
		if err != nil {
			panic(err)
		}

		remoteAddr, err := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
		if err != nil {
			panic(err)
		}

		_, err = conn.WriteToUDP(buildSearch(mx, target), remoteAddr)
		if err != nil {
			panic(err)
		}

		p := make([]byte, 2048)
		for {
			_, _, err = conn.ReadFromUDP(p)
			if err != nil {
				panic(err)
			}

			searchResults = append(searchResults, string(bytes.Trim(p, "\x00")))
		}
	}()

	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	for {
		select {
		case <-timer.C:
			return searchResults
		}
	}
}
