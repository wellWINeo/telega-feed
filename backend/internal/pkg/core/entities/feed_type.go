package entities

type FeedType = string

const (
	RSS  FeedType = "rss"
	RDF  FeedType = "rdf"
	Atom FeedType = "atom"
)

const DefaultFeedType = RSS
