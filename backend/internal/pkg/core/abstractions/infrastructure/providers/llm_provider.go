package abstractproviders

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
)

type LlmProvider interface {
	GenerateSummary(ctx context.Context, article *entities.Article) (string, error)
	GenerateDigest(ctx context.Context, articles []*entities.Article) (string, error)
}
