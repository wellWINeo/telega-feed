package abstractrepositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
)

type SummariesRepository interface {
	GetSummary(ctx context.Context, articleId entities.ArticleId) (*entities.Summary, error)
	AddSummary(ctx context.Context, articleId entities.ArticleId, summary *entities.Summary) error
}
