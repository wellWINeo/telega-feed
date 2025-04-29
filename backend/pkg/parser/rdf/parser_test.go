package rdf

import (
	"os"
	"testing"
)

func TestParseDateFromRSS_FSF(t *testing.T) {
	file, _ := os.Open("./fsf_news.xml")

	rdf, err := ParseRDF(file)
	if err != nil {
		t.Fatalf("ParseRSS failed: %v", err)
	}

	for i, item := range rdf.Items {
		if item.Title == "" {
			t.Fatalf("expected non-empty title, got %#v at %d", item.Title, i)
		}

		if item.Link == "" {
			t.Fatalf("expected non-empty link, got %#v at %d", item.Title, i)
		}
	}

	if len(rdf.Items) != 15 {
		t.Fatalf("expected 40 channels, got %d", len(rdf.Items))
	}
}
