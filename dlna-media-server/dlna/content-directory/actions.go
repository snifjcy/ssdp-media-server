package contentdirectory

import "encoding/xml"

type BrowseRequest struct {
	XMLName       xml.Name `xml:"Envelope"`
	S             string   `xml:"s,attr"`
	EncodingStyle string   `xml:"encodingStyle,attr"`
	Body          struct {
		U              string `xml:"u,attr"`
		ObjectID       string `xml:"ObjectID"`       // 0
		BrowseFlag     string `xml:"BrowseFlag"`     // BrowseDirectChildren
		Filter         string `xml:"Filter"`         // *
		StartingIndex  int    `xml:"StartingIndex"`  // 0
		RequestedCount int    `xml:"RequestedCount"` // 5000
		SortCriteria   string `xml:"SortCriteria"`
	} `xml:"Body>Browse"`
}

type BrowseResponse struct {
	XMLName       xml.Name `xml:"s:Envelope"`
	S             string   `xml:"xmlns:s,attr"`
	EncodingStyle string   `xml:"s:encodingStyle,attr"`
	Body          struct {
		U      string `xml:"xmlns:u,attr"`
		Result struct {
			Value string `xml:",innerxml"`
		} `xml:"Result"`
		NumberReturned string `xml:"NumberReturned"`
		TotalMatches   string `xml:"TotalMatches"`
		UpdateID       string `xml:"UpdateID"`
	} `xml:"s:Body>u:BrowseResponse"`
}

type BrowseResponseResult struct {
	XMLName   xml.Name                        `xml:"DIDL-Lite"`
	Xmlns     string                          `xml:"xmlns,attr"`
	Dc        string                          `xml:"xmlns:dc,attr"`
	Upnp      string                          `xml:"xmlns:upnp,attr"`
	Container []BrowseResponseResultContainer `xml:"container"`
}

type BrowseResponseResultContainer struct {
	ID         string `xml:"id,attr"`
	ParentID   string `xml:"parentID,attr"`
	Restricted string `xml:"restricted,attr"`
	ChildCount string `xml:"childCount,attr"`
	Title      string `xml:"dc:title"`
	Class      string `xml:"upnp:class"`
	Date       string `xml:"dc:date"`
	/*
		AlbumArtURI struct {
			Dlna      string `xml:"dlna,attr"`
			ProfileID string `xml:"profileID,attr"`
		} `xml:"albumArtURI"`
	*/
}

type Result struct {
	ID    string
	Title string
	URL   string
	ResultType
}

type ResultType int

const (
	ResultFolder ResultType = iota
	ResultImage
	ResultVideo
)
