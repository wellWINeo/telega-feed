package rdf

import "encoding/xml"

type RDF = rdf

type rdf struct {
	XMLName xml.Name `xml:"RDF"`
	Items   []item   `xml:"item"`
}

type item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Date        string `xml:"date"`
}
