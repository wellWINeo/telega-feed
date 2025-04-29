package providers

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
)

type StubLlmProvider struct {
}

func (s StubLlmProvider) GenerateSummary(ctx context.Context, article *entities.Article) (string, error) {
	return "Test summary", nil
}

func (s StubLlmProvider) GenerateDigest(ctx context.Context, articles []*entities.Article) (string, error) {
	return "Test digest", nil
}
