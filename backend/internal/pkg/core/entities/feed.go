package entities

type Feed struct {
	Articles []*Article `json:"articles,omitempty"`
	Digest   string     `json:"digest,omitempty"`
}
