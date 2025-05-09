package entities

import "github.com/google/uuid"

type FeedSourceId = uuid.UUID

type FeedSource struct {
	Id       FeedSourceId `json:"id" sql:"id"`
	Name     string       `json:"name" sql:"name"`
	FeedUrl  string       `json:"feed_url" sql:"feed_url"`
	Type     FeedType     `json:"type" sql:"type"`
	Disabled bool         `json:"disabled" sql:"disabled"`
}
