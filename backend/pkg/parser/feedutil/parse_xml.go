package feedutil

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"io"
)

func ParseXML[T any](input io.Reader) (*T, error) {
	var result T

	decoder := xml.NewDecoder(input)
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
