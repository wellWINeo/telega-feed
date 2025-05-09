package entities

import (
	"github.com/google/uuid"
	"time"
)

type ArticleId = uuid.UUID

type Article struct {
	Id          ArticleId    `json:"id,omitempty" sql:"id"`
	SourceId    FeedSourceId `json:"sourceId,omitempty" sql:"source_id"`
	AddedAt     time.Time    `json:"added_at,omitempty" sql:"added_at"`
	PublishedAt time.Time    `json:"published_at,omitempty" sql:"published_at"`
	Title       string       `json:"title,omitempty" sql:"title"`
	Text        string       `json:"text,omitempty" sql:"text"`
	Url         string       `json:"url,omitempty" sql:"url"`
	PreviewUrl  string       `json:"preview_url,omitempty" sql:"preview_url"`
	Starred     bool         `json:"starred,omitempty" sql:"starred"`
	Read        bool         `json:"read,omitempty" sql:"read"`
}
