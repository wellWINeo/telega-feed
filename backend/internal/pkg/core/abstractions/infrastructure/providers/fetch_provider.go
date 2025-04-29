package abstractproviders

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
	"io"
)

type FetchProvider interface {
	CheckType(reader io.Reader) (bool, error)
	FetchArticles(ctx context.Context, feedSource *entities.FeedSource) ([]*entities.Article, error)
}
