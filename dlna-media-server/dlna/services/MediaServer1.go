package services

import (
	"encoding/xml"
	"fmt"
)

type Device root
type ContentDirectoryService Service
type ConnectionManagerService Service
type AVTransportService Service

type root struct {
	XMLName     xml.Name `xml:"root"`
	Xmlns       string   `xml:"xmlns,attr"` // must be `urn:schemas-upnp-org:device-1-0`
	SpecVersion struct {
		Major int `xml:"major"` // must be 1
		Minor int `xml:"minor"` // must be 0
	} `xml:"specVersion"`
	URLBase string            `xml:"URLBase"` // base URL for all relative URLs
	Device  deviceDescription `xml:"device"`
}

type deviceDescription struct {
	DeviceType       string    `xml:"deviceType"`       // must be `urn:schemas-upnp-org:device:MediaServer:1`
	FriendlyName     string    `xml:"friendlyName"`     // short user-friendly title
	Manufacturer     string    `xml:"manufacturer"`     // manufacturer name
	ManufacturerURL  string    `xml:"manufacturerURL"`  // URL to manufacturer site
	ModelDescription string    `xml:"modelDescription"` // long user-friendly title
	ModelName        string    `xml:"modelName"`        // model name
	ModelNumber      string    `xml:"modelNumber"`      // model number
	ModelURL         string    `xml:"modelURL"`         // URL to model site
	SerialNumber     string    `xml:"serialNumber"`     // manufacturer's serial number
	UDN              string    `xml:"UDN"`              // must be `uuid:%s`
	UPC              string    `xml:"UPC"`              // Universal Product Code
	ServiceList      []Service `xml:"serviceList>service"`
	PresentationURL  string    `xml:"presentationURL"` // URL for presentation
}

type Service struct {
	ServiceType string `xml:"serviceType"` // must be `urn:schemas-upnp-org:service:%s:1` depending on the service
	ServiceId   string `xml:"serviceId"`   // must be `urn:upnp-org:serviceId:%s` depending on the service
	SCPDURL     string `xml:"SCPDURL"`     // path of the URL to service description
	ControlURL  string `xml:"controlURL"`  // path of the URL for control
	EventSubURL string `xml:"eventSubURL"` // path of the URL for eventing
}

func NewDevice(name string, uuid string) Device {
	r := Device{
		Xmlns: "urn:schemas-upnp-org:device-1-0",
		SpecVersion: struct {
			Major int `xml:"major"`
			Minor int `xml:"minor"`
		}{
			Major: 1,
			Minor: 0,
		},
		URLBase: "http://192.168.1.95:8090",
	}

	r.Device = deviceDescription{
		DeviceType:       "urn:schemas-upnp-org:device:MediaServer:1",
		FriendlyName:     name,
		ModelName:        "dlna-media-server",
		Manufacturer:     "Poup2804",
		ManufacturerURL:  "https://github.com/Poup2804",
		ModelDescription: name,
		ModelNumber:      "0.1",
		ModelURL:         "https://github.com/Poup2804/dlna-media-server",
		SerialNumber:     "2804",
		UDN:              fmt.Sprintf("uuid:%s", uuid),
	}

	return r
}

func (d *Device) SetContentDirectoryService(s Service) {
	const ContentDirectoryId = "urn:upnp-org:serviceId:ContentDirectory"
	const ContentDirectoryType = "urn:schemas-upnp-org:service:ContentDirectory:1"

	serviceFound := false

	for _, service := range d.Device.ServiceList {
		if service.ServiceId == ContentDirectoryId {
			serviceFound = true
			service.SCPDURL = s.SCPDURL
			service.ControlURL = s.ControlURL
			service.EventSubURL = s.EventSubURL
			break
		}
	}

	if !serviceFound {
		s.ServiceId = ContentDirectoryId
		s.ServiceType = ContentDirectoryType
		d.Device.ServiceList = append(d.Device.ServiceList, s)
	}
}
