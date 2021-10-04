package main

import (
	"dlna-media-server/dlna/content-directory"
	"dlna-media-server/dlna/services"
	"encoding/xml"
	"fmt"
	"github.com/koron/go-ssdp"
	"log"
	"net/http"
	"time"
)

func main() {
	uuid := "d726b503-d259-4633-bb8c-b7d37960adf1"
	device := services.NewDevice("My media server", uuid)
	device.SetContentDirectoryService(services.Service{
		SCPDURL:     fmt.Sprintf("/%s/content-directory/scpd.xml", uuid),
		ControlURL:  fmt.Sprintf("/%s/content-directory/control", uuid),
		EventSubURL: fmt.Sprintf("/%s/content-directory/subscribe", uuid),
	})
	r, e := xml.Marshal(device)
	if e != nil {
		log.Panicln(e)
	}

	http.HandleFunc(fmt.Sprintf("/%s/desc.xml", uuid), func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "text/xml; charset=\"UTF-8\"")
		writer.Write(r)
	})
	http.HandleFunc(fmt.Sprintf("/%s/content-directory/control", uuid), func(writer http.ResponseWriter, request *http.Request) {
		var BR contentdirectory.BrowseRequest
		err := xml.NewDecoder(request.Body).Decode(&BR)
		if err != nil {
			log.Println("Porco diooo", err)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte{})
			return
		}
		res, err := contentdirectory.BuildBrowseResponse([]contentdirectory.Result{
			contentdirectory.Result{Title: "Titolo 1", ID: "1001", ResultType: contentdirectory.ResultFolder},
			contentdirectory.Result{Title: "Titolo 2", ID: "1002", ResultType: contentdirectory.ResultFolder},
		})
		if err != nil {
			log.Println("Porco diooo2", err)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte{})
			return
		}
		writer.Header().Add("Content-Type", "text/xml; charset=\"UTF-8\"")
		writer.Write([]byte("<?xml version='1.0' encoding='utf-8'?>\n"))
		writer.Write(res)
	})
	go startSSDP("192.168.1.95", uuid)
	http.ListenAndServe("192.168.1.95:8090", nil)
}

func startSSDP(ip string, uuid string) {
	ad, err := ssdp.Advertise(
		"urn:schemas-upnp-org:device:MediaServer:1",                             // send as "ST"
		fmt.Sprintf("uuid:%s::urn:schemas-upnp-org:device:MediaServer:1", uuid), // send as "USN"
		fmt.Sprintf("http://%s:8090/%s/desc.xml", ip, uuid),                     // send as "LOCATION"
		"go-ssdp sample", // send as "SERVER"
		1800)             // send as "maxAge" in "CACHE-CONTROL"
	if err != nil {
		panic(err)
	}
	aliveTick := time.Tick(300 * time.Second)

	for {
		select {
		case <-aliveTick:
			ad.Alive()
		}
	}
}
