package services

import (
	"TelegaFeed/internal/pkg/core/abstractions/infrastructure/providers"
	"TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"
	"TelegaFeed/internal/pkg/core/entities"
	"TelegaFeed/pkg/mytime"
	"context"
	"log"
)

type LlmService struct {
	digestsRepository   abstractrepositories.DigestsRepository
	summariesRepository abstractrepositories.SummariesRepository
	feedRepository      abstractrepositories.FeedRepository
	llmProvider         abstractproviders.LlmProvider
}

func NewLlmService(
	digestsRepository abstractrepositories.DigestsRepository,
	summariesRepository abstractrepositories.SummariesRepository,
	feedRepository abstractrepositories.FeedRepository,
	llmProvider abstractproviders.LlmProvider,
) *LlmService {
	return &LlmService{
		digestsRepository:   digestsRepository,
		summariesRepository: summariesRepository,
		feedRepository:      feedRepository,
		llmProvider:         llmProvider,
	}
}

func (l *LlmService) GetDailyDigest(ctx context.Context, userId entities.UserId) (*string, error) {
	ok, existedDigest, err := l.digestsRepository.FindLatestDigestForToday(ctx, userId)
	if err != nil {
		return nil, err
	}

	if ok {
		return &existedDigest, nil
	}

	return l.GenerateDailyDigest(ctx, userId)
}

func (l *LlmService) GenerateDailyDigest(ctx context.Context, userId entities.UserId) (*string, error) {
	todayArticles, err := l.feedRepository.GetTodayArticles(ctx, userId)
	if err != nil {
		return nil, err
	}

	generatedDigest, err := l.llmProvider.GenerateDigest(ctx, todayArticles)
	if err != nil {
		return nil, err
	}

	err = l.digestsRepository.AddDigest(ctx, userId, generatedDigest)
	if err != nil {
		log.Printf("Failed to save digest to db: %v", err)
	}

	return &generatedDigest, nil
}

func (l *LlmService) GetArticleSummary(
	ctx context.Context,
	userId entities.UserId,
	articleId entities.ArticleId,
) (*entities.Summary, error) {
	existedSummary, err := l.summariesRepository.GetSummary(ctx, articleId)
	if err != nil {
		return nil, err
	}

	if existedSummary != nil {
		return existedSummary, nil
	}

	article, err := l.feedRepository.GetArticleById(ctx, userId, articleId)
	if err != nil {
		return nil, err
	}

	summaryText, err := l.llmProvider.GenerateSummary(ctx, article)
	if err != nil {
		return nil, err
	}

	generatedSummary := entities.Summary{
		GeneratedAt: mytime.NowUTC(),
		Text:        summaryText,
	}

	err = l.summariesRepository.AddSummary(ctx, articleId, &generatedSummary)
	if err != nil {
		log.Printf("Failed to save summary to db: %v", err)
	}

	return &generatedSummary, nil
}
