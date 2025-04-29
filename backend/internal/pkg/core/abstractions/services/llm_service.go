package abstractservices

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
)

type LlmService interface {
	GetDailyDigest(ctx context.Context, userId entities.UserId) (*string, error)
	GenerateDailyDigest(ctx context.Context, userId entities.UserId) (*string, error)
	GetArticleSummary(ctx context.Context, userId entities.UserId, articleId entities.ArticleId) (*entities.Summary, error)
}
