package rss

import (
	"TelegaFeed/pkg/parser/feedutil"
	"fmt"
	"io"
	"strings"
)

func ParseRSS(input io.Reader) (*RSS, error) {
	rss, err := feedutil.ParseXML[RSS](input)
	if err != nil {
		return nil, fmt.Errorf("error decoding RSS: %w", err)
	}

	return rss, nil
}

func GetPreviewUrlFromRSS(enclosures []enclosure) string {
	for _, enc := range enclosures {
		if strings.HasPrefix(enc.Type, "image/") {
			return enc.URL
		}
	}

	if len(enclosures) > 0 {
		return enclosures[0].URL
	}

	return ""
}
