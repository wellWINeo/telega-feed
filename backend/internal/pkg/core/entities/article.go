package entities

import "time"

type ArticleId = string

type Article struct {
	Id          ArticleId `json:"id,omitempty"`
	AddedAt     time.Time `json:"added_at,omitempty"`
	PublishedAt time.Time `json:"published_at,omitempty"`
	Title       string    `json:"title,omitempty"`
	Text        string    `json:"text,omitempty"`
	Url         string    `json:"url,omitempty"`
	PreviewUrl  string    `json:"preview_url,omitempty"`
	Starred     bool      `json:"starred,omitempty"`
	Read        bool      `json:"read,omitempty"`
}
