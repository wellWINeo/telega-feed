package atom

import (
	"TelegaFeed/pkg/parser/feedutil"
	"fmt"
	"io"
	"strings"
)

func ParseAtom(input io.Reader) (*Feed, error) {
	feed, err := feedutil.ParseXML[Feed](input)

	if err != nil {
		return nil, fmt.Errorf("error decoding Atom feed: %w", err)
	}

	return feed, nil
}

func GetLink(links []link) string {
	for _, link := range links {
		if link.Rel == "alternate" || link.Rel == "" {
			return link.Href
		}
	}

	return ""
}

func GetPreviewLink(links []link) string {
	for _, link := range links {
		if link.Rel == "enclosure" && strings.HasPrefix(link.Href, "image/") {
			return link.Href
		}
	}

	return ""
}
