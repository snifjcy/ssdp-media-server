package main

import (
	"errors"
	"fmt"
	"github.com/koron/go-ssdp"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	localIp, err := getIP()
	if err != nil {
		log.Panicf("Couldn't get IP:\n%+v\n", err)
	}

	go func(ip net.IP) {
		ad, err := ssdp.Advertise(
			"urn:schemas-upnp-org:device:MediaServer:1",                                                       // send as "ST"
			"unique:id::urn:schemas-upnp-org:device:MediaServer:1",                                            // send as "USN"
			fmt.Sprintf("http://%s:8080/8077b5ab-d117-4893-aafd-552529ddc263/description-1.xml", ip.String()), // send as "LOCATION"
			"go-ssdp sample", // send as "SERVER"
			1800)             // send as "maxAge" in "CACHE-CONTROL"
		if err != nil {
			panic(err)
		}
		defer func() {
			ad.Bye()
			fmt.Println("Byebyed")
			ad.Close()
			fmt.Println("Closed")
		}()
		aliveTick := time.Tick(300 * time.Second)
		for {
			select {
			case <-aliveTick:
				ad.Alive()
			}
		}
	}(localIp)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

func getIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return net.IP{}, err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}
	return net.IP{}, errors.New("No ip found")
}
