package MediaServer1

import (
	"encoding/xml"
	"fmt"
)

type Device rootDescription

type rootDescription struct {
	XMLName     xml.Name          `xml:"root"`
	Xmlns       string            `xml:"xmlns,attr"`
	SpecVersion specVersion       `xml:"specVersion"`
	device      deviceDescription `xml:"device"`
}

type specVersion struct {
	Major int `xml:"major"`
	Minor int `xml:"minor"`
}

type deviceDescription struct {
	DeviceType       string `xml:"deviceType"`
	FriendlyName     string `xml:"friendlyName"`
	ModelName        string `xml:"modelName"`
	Manufacturer     string `xml:"manufacturer"`
	ManufacturerURL  string `xml:"manufacturerURL"`
	ModelDescription string `xml:"modelDescription"`
	ModelNumber      string `xml:"modelNumber"`
	ModelURL         string `xml:"modelURL"`
	SerialNumber     string `xml:"serialNumber"`
	UDN              string `xml:"UDN"`
	UPC              string `xml:"UPC"`
	ServiceList      struct {
		Service []serviceDescription `xml:"service"`
	} `xml:"serviceList"`
	PresentationURL string   `xml:"presentationURL"`
	XDLNADOC        []string `xml:"X_DLNADOC"`
	XDLNACAP        string   `xml:"X_DLNACAP"`
}

type serviceDescription struct {
	ServiceType string `xml:"serviceType"`
	ServiceId   string `xml:"serviceId"`
	SCPDURL     string `xml:"SCPDURL"`
	ControlURL  string `xml:"controlURL"`
	EventSubURL string `xml:"eventSubURL"`
}

func NewDevice(name string, uuid string) Device {
	r := &rootDescription{
		Xmlns: "urn:schemas-upnp-org:device-1-0",
		SpecVersion: specVersion{
			Major: 1,
			Minor: 0,
		},
	}

	d := &deviceDescription{
		DeviceType:       "urn:schemas-upnp-org:device:MediaServer:1",
		FriendlyName:     name,
		ModelName:        "go-drive-test",
		Manufacturer:     "Poup2804",
		ManufacturerURL:  "https://github.com/Poup2804",
		ModelDescription: name,
		ModelNumber:      "0.1",
		ModelURL:         "https://github.com/Poup2804/go-drive-test",
		SerialNumber:     "2804",
		UDN:              fmt.Sprintf("uuid:%s", uuid),
		UPC:              "",
		ServiceList: struct {
			Service []serviceDescription `xml:"service"`
		}{},
		PresentationURL: "",
		XDLNADOC:        nil,
		XDLNACAP:        "",
	}
}

/*
<?xml version='1.0' encoding='utf-8'?>
<root
	xmlns=>
	<specVersion>
		<major>1</major>
		<minor>0</minor>
	</specVersion>
	<device>
		<deviceType>urn:schemas-upnp-org:device:MediaServer:1</deviceType>
		<friendlyName>my media</friendlyName>
		<modelName>Cohen3</modelName>
		<manufacturer>beebits.net</manufacturer>
		<manufacturerURL>https://github.com/opacam/Cohen3</manufacturerURL>
		<modelDescription>Cohen3</modelDescription>
		<modelNumber>0.9.3</modelNumber>
		<modelURL>https://github.com/opacam/Cohen3</modelURL>
		<serialNumber>0000001</serialNumber>
		<UDN>uuid:8077b5ab-d117-4893-aafd-552529ddc263</UDN>
		<UPC></UPC>
		<serviceList>
			<service>
				<serviceType>urn:schemas-upnp-org:service:ConnectionManager:1</serviceType>
				<serviceId>urn:upnp-org:serviceId:ConnectionManager</serviceId>
				<SCPDURL>/8077b5ab-d117-4893-aafd-552529ddc263/ConnectionManager/scpd.xml</SCPDURL>
				<controlURL>/8077b5ab-d117-4893-aafd-552529ddc263/ConnectionManager/control</controlURL>
				<eventSubURL>/8077b5ab-d117-4893-aafd-552529ddc263/ConnectionManager/subscribe</eventSubURL>
			</service>
			<service>
				<serviceType>urn:schemas-upnp-org:service:ContentDirectory:1</serviceType>
				<serviceId>urn:upnp-org:serviceId:ContentDirectory</serviceId>
				<SCPDURL>/8077b5ab-d117-4893-aafd-552529ddc263/ContentDirectory/scpd.xml</SCPDURL>
				<controlURL>/8077b5ab-d117-4893-aafd-552529ddc263/ContentDirectory/control</controlURL>
				<eventSubURL>/8077b5ab-d117-4893-aafd-552529ddc263/ContentDirectory/subscribe</eventSubURL>
			</service>
		</serviceList>
		<presentationURL>/8077b5ab-d117-4893-aafd-552529ddc263</presentationURL>
		<X_DLNADOC>DMS-1.50</X_DLNADOC>
		<X_DLNADOC>M-DMS-1.50</X_DLNADOC>
		<X_DLNACAP>av-upload,image-upload,audio-upload</X_DLNACAP>
	</device>
</root>
*/
