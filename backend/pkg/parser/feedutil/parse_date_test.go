package feedutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getFixedTimeProvider(fixedTime time.Time) func() time.Time {
	return func() time.Time {
		return fixedTime
	}
}

func TestParseFeedDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "RFC1123Z format",
			input:    "Mon, 02 Jan 2006 15:04:05 -0700",
			expected: time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60)),
		},
		{
			name:     "RFC3339 format",
			input:    "2006-01-02T15:04:05Z",
			expected: time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
		},
		{
			name:     "Simple date format",
			input:    "2006-01-02",
			expected: time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "RDF-specific date format",
			input:    "2006-01-02T15:04:05-07:00",
			expected: time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60)),
		},
		{
			name:     "Invalid date",
			input:    "not-a-date",
			expected: time.Date(2025, 4, 14, 21, 0, 0, 0, time.UTC),
		},
		{
			name:     "Empty date",
			input:    "",
			expected: time.Date(2025, 4, 14, 21, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixedTime := getFixedTimeProvider(tt.expected)
			got := parseFeedDate(tt.input, fixedTime)

			assert.Equal(t, tt.expected, got)
		})
	}
}
