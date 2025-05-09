package abstractservices

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
)

type FeedService interface {
	GetFeed(ctx context.Context, userId entities.UserId) (*entities.Feed, error)
	UpdateArticle(ctx context.Context, userId entities.UserId, articleId entities.ArticleId, patch *entities.ArticlePatch) error
	UpdateFeed(ctx context.Context) error
}
