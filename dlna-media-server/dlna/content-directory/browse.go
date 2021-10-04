package contentdirectory

import (
	"encoding/xml"
	"fmt"
	"log"
)

const browseResponseFormat = `<?xml version="1.0" encoding="utf-8" standalone="yes"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body>%s</s:Body></s:Envelope>`
const didlLiteFormat = `<?xml version="1.0" encoding="utf-8" standalone="yes"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body>%s</s:Body></s:Envelope>`
const abcFormat = `<u:%[1]sResponse xmlns:u="%[2]s">%[3]s</u:%[1]sResponse>`

type Container struct {
	Object
	XMLName    xml.Name `xml:"container"`
	ChildCount int      `xml:"childCount,attr"`
}

type Object struct {
	ID         string      `xml:"id,attr"`
	ParentID   string      `xml:"parentID,attr"`
	Restricted int         `xml:"restricted,attr"`
	Class      ObjectClass `xml:"upnp:class"`
	Icon       string      `xml:"upnp:icon,omitempty"`
	Title      string      `xml:"dc:title"`
	//Date        Timestamp  `xml:"dc:date"`
	//Artist      string     `xml:"upnp:artist,omitempty"`
	//Album       string     `xml:"upnp:album,omitempty"`
	//Genre       string     `xml:"upnp:genre,omitempty"`
	//AlbumArtURI string     `xml:"upnp:albumArtURI,omitempty"`
	Searchable int `xml:"searchable,attr"`
}

type ObjectClass string

const (
	ObjectClassFolder ObjectClass = "object.container.storageFolder"
	ObjectClassVideo  ObjectClass = "object.item.videoItem"
)

func NewObject(isFolder bool, id string, parentID string, title string) (interface{}, error) {
	obj := Object{
		ID:       id,
		ParentID: parentID,
		Title:    title,
	}

	if isFolder {
		obj.Class = ObjectClassFolder
		return Container{
			Object:     obj,
			ChildCount: 1,
		}, nil
	}

	return obj, nil
}

func BuildBrowseResponseV2(results []Result) ([]byte, error) {
	var listOfResults []interface{}
	for _, result := range results {
		obj, err := NewObject(true, result.ID, "1000", result.Title)
		if err != nil {
			log.Println("Errore nella conversione in buildbrowseresponsev2")
		}
		listOfResults = append(listOfResults, obj)
	}

	result, err := xml.Marshal(listOfResults)
	if err != nil {
		return nil, err
	}
	test := map[string]string{
		"TotalMatches":   fmt.Sprint(len(result)),
		"NumberReturned": fmt.Sprint(len(result)),
		"Result":         fmt.Sprintf(didlLiteFormat, string(result)),
		"UpdateID":       "1",
	}

	result2, err := xml.Marshal(test)
	if err != nil {
		return nil, err
	}

	test2 := fmt.Sprintf(abcFormat, "browse", "urn:schemas-upnp-org:service:content-directory:1", result2)
	test3 := fmt.Sprintf(browseResponseFormat, test2)

	return []byte(test3), nil
}
