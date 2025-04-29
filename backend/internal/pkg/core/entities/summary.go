package entities

import "time"

type SummaryId = string

type Summary struct {
	Id          SummaryId `json:"id,omitempty"`
	GeneratedAt time.Time `json:"generated_at,omitempty"`
	Text        string    `json:"text,omitempty"`
}
