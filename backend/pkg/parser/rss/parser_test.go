package rss

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func TestParseRSS_ValidInput(t *testing.T) {
	input := strings.NewReader(`
		<rss>
			<channel>
				<item>
					<title>Test Article</title>
					<link>http://example.com/article</link>
					<description>This is a test article.</description>
					<pubDate>Mon, 02 Jan 2006 15:04:05 MST</pubDate>
					<enclosure url="http://example.com/image.jpg" type="image/jpeg" />
				</item>
			</channel>
		</rss>
	`)

	rss, err := ParseRSS(input)
	if err != nil {
		t.Fatalf("ParseRSS failed: %v", err)
	}

	if len(rss.Channel.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(rss.Channel.Items))
	}

	item := rss.Channel.Items[0]
	if item.Title != "Test Article" {
		t.Errorf("expected title 'Test Article', got '%s'", item.Title)
	}
	if item.Link != "http://example.com/article" {
		t.Errorf("expected link 'http://example.com/article', got '%s'", item.Link)
	}
	if item.Description != "This is a test article." {
		t.Errorf("expected description 'This is a test article.', got '%s'", item.Description)
	}
	if item.PubDate != "Mon, 02 Jan 2006 15:04:05 MST" {
		t.Errorf("expected pubDate 'Mon, 02 Jan 2006 15:04:05 MST', got '%s'", item.PubDate)
	}
	if len(item.Enclosures) != 1 {
		t.Fatalf("expected 1 enclosure, got %d", len(item.Enclosures))
	}
	if item.Enclosures[0].URL != "http://example.com/image.jpg" {
		t.Errorf("expected enclosure URL 'http://example.com/image.jpg', got '%s'", item.Enclosures[0].URL)
	}
	if item.Enclosures[0].Type != "image/jpeg" {
		t.Errorf("expected enclosure type 'image/jpeg', got '%s'", item.Enclosures[0].Type)
	}
}

func TestParseRSS_InvalidInput(t *testing.T) {
	input := strings.NewReader(`<invalid><xml></xml></invalid>`)

	_, err := ParseRSS(input)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "error decoding RSS") {
		t.Errorf("expected error to contain 'error decoding RSS', got '%v'", err)
	}
}

func TestParseRSS_EmptyInput(t *testing.T) {
	input := strings.NewReader("")

	_, err := ParseRSS(input)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, io.EOF) {
		t.Errorf("expected EOF error, got '%v'", err)
	}
}

func TestGetPreviewUrlFromRSS_ImageEnclosure(t *testing.T) {
	enclosures := []enclosure{
		{URL: "http://example.com/image.jpg", Type: "image/jpeg"},
		{URL: "http://example.com/video.mp4", Type: "video/mp4"},
	}

	url := GetPreviewUrlFromRSS(enclosures)
	if url != "http://example.com/image.jpg" {
		t.Errorf("expected 'http://example.com/image.jpg', got '%s'", url)
	}
}

func TestGetPreviewUrlFromRSS_NonImageEnclosure(t *testing.T) {
	enclosures := []enclosure{
		{URL: "http://example.com/video.mp4", Type: "video/mp4"},
		{URL: "http://example.com/audio.mp3", Type: "audio/mp3"},
	}

	url := GetPreviewUrlFromRSS(enclosures)
	if url != "http://example.com/video.mp4" {
		t.Errorf("expected 'http://example.com/video.mp4', got '%s'", url)
	}
}

func TestGetPreviewUrlFromRSS_NoEnclosures(t *testing.T) {
	var enclosures []enclosure

	url := GetPreviewUrlFromRSS(enclosures)
	if url != "" {
		t.Errorf("expected empty string, got '%s'", url)
	}
}

func TestParseDateFromRSS_OpenNet(t *testing.T) {
	file, _ := os.Open("./opennews_all.rss")

	rss, err := ParseRSS(file)
	if err != nil {
		t.Fatalf("ParseRSS failed: %v", err)
	}

	if len(rss.Channel.Items) != 40 {
		t.Fatalf("expected 40 channels, got %d", len(rss.Channel.Items))
	}
}
