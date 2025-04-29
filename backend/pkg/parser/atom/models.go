package atom

import "encoding/xml"

type Feed = feed

type feed struct {
	XMLName xml.Name `xml:"feed"`
	Title   string   `xml:"title"`
	Entries []entry  `xml:"entry"`
}

type entry struct {
	Title   string   `xml:"title"`
	Links   []link   `xml:"link"`
	Updated string   `xml:"updated"`
	Content *content `xml:"content"`
}

type link struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type content struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}
