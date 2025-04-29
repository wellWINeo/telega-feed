package rdf

import (
	"TelegaFeed/pkg/parser/feedutil"
	"io"
)

func ParseRDF(input io.Reader) (*RDF, error) {
	rdf, err := feedutil.ParseXML[RDF](input)

	if err != nil {
		return nil, err
	}

	return rdf, nil
}
