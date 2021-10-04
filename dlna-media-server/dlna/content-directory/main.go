package contentdirectory

import (
	"encoding/xml"
	"strings"

	//"html"
	"log"
)

func ParseBrowseRequest(body []byte) (BrowseRequest, error) {
	var BR BrowseRequest
	err := xml.Unmarshal(body, &BR)
	if err != nil {
		log.Println(err)
		return BrowseRequest{}, err
	}
	return BR, nil
}

func BuildBrowseResponse(results []Result) ([]byte, error) {
	BR := BrowseResponse{
		S:             "http://schemas.xmlsoap.org/soap/envelope/",
		EncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		Body: struct {
			U      string `xml:"xmlns:u,attr"`
			Result struct {
				Value string `xml:",innerxml"`
			} `xml:"Result"`
			NumberReturned string `xml:"NumberReturned"`
			TotalMatches   string `xml:"TotalMatches"`
			UpdateID       string `xml:"UpdateID"`
		}{
			U:              "urn:schemas-upnp-org:service:content-directory:1",
			NumberReturned: "2",
			TotalMatches:   "2",
			UpdateID:       "2",
		},
	}

	BRR := BrowseResponseResult{
		Xmlns: "urn:schemas-upnp-org:metadata-1-0/DIDL-Lite/",
		Dc:    "http://purl.org/dc/elements/1.1/",
		Upnp:  "urn:schemas-upnp-org:metadata-1-0/upnp/",
	}

	for _, result := range results {
		BRRC := BrowseResponseResultContainer{
			ID:         result.ID,
			ParentID:   "1000",
			Title:      result.Title,
			Date:       "1999-09-19T04:12:00+02:00",
			ChildCount: "1",
			Restricted: "0",
		}
		switch result.ResultType {
		case ResultFolder:
			BRRC.Class = "object.container"
			break
		case ResultImage:
			BRRC.Class = "object.image"
			break
		case ResultVideo:
			BRRC.Class = "object.video"
			break
		}

		BRR.Container = append(BRR.Container, BRRC)
	}

	unescapedBRR, err := xml.Marshal(BRR)
	e := strings.ReplaceAll(string(unescapedBRR), "<", "&lt;")
	e = strings.ReplaceAll(e, ">", "&gt;")
	//escapedBRR := html.EscapeString(string(unescapedBRR))
	BR.Body.Result.Value = e

	marshaledBR, err := xml.Marshal(BR)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return marshaledBR, nil
}
