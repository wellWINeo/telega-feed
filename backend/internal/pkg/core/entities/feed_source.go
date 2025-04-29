package entities

type FeedSourceId = string

type FeedSource struct {
	Id       FeedSourceId `json:"id"`
	Name     string       `json:"name"`
	FeedUrl  string       `json:"feed_url"`
	Type     FeedType     `json:"type"`
	Disabled bool         `json:"disabled"`
}
