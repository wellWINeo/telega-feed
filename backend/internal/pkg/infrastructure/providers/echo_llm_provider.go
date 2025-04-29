package providers

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
)

type EchoLlmProvider struct {
}

func NewEchoLlmProvider() *EchoLlmProvider {
	return &EchoLlmProvider{}
}

func (e EchoLlmProvider) GenerateSummary(ctx context.Context, article *entities.Article) (string, error) {
	return reverse(article.Text), nil
}

func (e EchoLlmProvider) GenerateDigest(ctx context.Context, articles []*entities.Article) (string, error) {
	return "sample digest", nil
}

func reverse(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
