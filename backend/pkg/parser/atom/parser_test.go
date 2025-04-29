package atom

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestParseAtom_ValidInput(t *testing.T) {
	input := strings.NewReader(`
		<feed xmlns="http://www.w3.org/2005/Atom">
			<title>Example Feed</title>
			<entry>
				<title>Atom-Powered Robots Run Amok</title>
				<link href="http://example.com/1" rel="alternate" type="text/html" />
				<updated>2003-12-13T18:30:02Z</updated>
				<content type="text">Some text content</content>
			</entry>
		</feed>
	`)

	feed, err := ParseAtom(input)
	if err != nil {
		t.Fatalf("ParseAtom failed: %v", err)
	}

	if feed.Title != "Example Feed" {
		t.Errorf("expected title 'Example Feed', got '%s'", feed.Title)
	}

	if len(feed.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(feed.Entries))
	}

	entry := feed.Entries[0]
	if entry.Title != "Atom-Powered Robots Run Amok" {
		t.Errorf("expected entry title 'Atom-Powered Robots Run Amok', got '%s'", entry.Title)
	}

	if entry.Updated != "2003-12-13T18:30:02Z" {
		t.Errorf("expected updated '2003-12-13T18:30:02Z', got '%s'", entry.Updated)
	}

	if entry.Content == nil || entry.Content.Value != "Some text content" {
		t.Errorf("expected content 'Some text content', got '%v'", entry.Content)
	}
}

func TestParseAtom_InvalidInput(t *testing.T) {
	input := strings.NewReader(`<invalid><xml></xml></invalid>`)

	_, err := ParseAtom(input)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "error decoding Atom feed") {
		t.Errorf("expected error to contain 'error decoding Atom feed', got '%v'", err)
	}
}

func TestParseAtom_EmptyInput(t *testing.T) {
	input := strings.NewReader("")

	_, err := ParseAtom(input)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, io.EOF) {
		t.Errorf("expected EOF error, got '%v'", err)
	}
}

func TestGetAlternateLink_Present(t *testing.T) {
	links := []link{
		{Href: "http://example.com/1", Rel: "alternate", Type: "text/html"},
		{Href: "http://example.com/2", Rel: "enclosure", Type: "image/jpeg"},
	}

	url := GetLink(links)
	if url != "http://example.com/1" {
		t.Errorf("expected 'http://example.com/1', got '%s'", url)
	}
}

func TestGetAlternateLink_Absent(t *testing.T) {
	links := []link{
		{Href: "http://example.com/2", Rel: "enclosure", Type: "image/jpeg"},
	}

	url := GetLink(links)
	if url != "" {
		t.Errorf("expected empty string, got '%s'", url)
	}
}

func TestGetAlternateLink_Empty(t *testing.T) {
	var links []link

	url := GetLink(links)
	if url != "" {
		t.Errorf("expected empty string, got '%s'", url)
	}
}
