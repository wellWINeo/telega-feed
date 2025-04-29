package feedutil

import (
	"strings"
	"time"
)

var dateLayouts = []string{
	time.RFC1123Z,                  // RSS standard format
	time.RFC822Z,                   // Alternate RSS format
	time.RFC3339,                   // Atom/RDF format
	"2006-01-02",                   // Simple date
	"02 Jan 2006",                  // Alternative date format
	time.RFC1123,                   // With timezone name
	"Mon, 2 Jan 2006 15:04:05 MST", // Common RSS variant
}

func ParseFeedDate(dateStr string) time.Time {
	return parseFeedDate(dateStr, time.Now)
}

func parseFeedDate(dateStr string, timeProvider func() time.Time) time.Time {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return timeProvider().UTC()
	}

	// Try each layout until one succeeds
	for _, layout := range dateLayouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return t
		}
	}

	return timeProvider().UTC()
}
