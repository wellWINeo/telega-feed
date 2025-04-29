package abstractservices

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
)

type FetchService interface {
	DetectType(feedSource *entities.FeedSource) entities.FeedType
	FetchArticles(ctx context.Context, feedSources []*entities.FeedSource) []*entities.Article
}
