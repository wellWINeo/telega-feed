package services

import (
	"TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"
	"TelegaFeed/internal/pkg/core/abstractions/services"
	"TelegaFeed/internal/pkg/core/entities"
	"context"
)

type FeedService struct {
	llmService           abstractservices.LlmService
	fetchService         abstractservices.FetchService
	feedRepository       abstractrepositories.FeedRepository
	feedSourceRepository abstractrepositories.FeedSourceRepository
	usersRepository      abstractrepositories.UsersRepository
}

func NewFeedService(
	llmService abstractservices.LlmService,
	fetchService abstractservices.FetchService,
	feedRepository abstractrepositories.FeedRepository,
	feedSourceRepository abstractrepositories.FeedSourceRepository,
	usersRepository abstractrepositories.UsersRepository,
) *FeedService {
	return &FeedService{
		llmService:           llmService,
		fetchService:         fetchService,
		feedRepository:       feedRepository,
		feedSourceRepository: feedSourceRepository,
		usersRepository:      usersRepository,
	}
}

func (f *FeedService) GetFeed(ctx context.Context, userId entities.UserId) (*entities.Feed, error) {
	articles, err := f.feedRepository.GetFeedByUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	digest, err := f.llmService.GetDailyDigest(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &entities.Feed{
		Articles: articles,
		Digest:   *digest,
	}, nil
}

func (f *FeedService) UpdateArticle(
	ctx context.Context,
	userId entities.UserId,
	articleId entities.ArticleId,
	patch *entities.ArticlePatch,
) error {
	_, err := f.feedRepository.UpdateArticle(ctx, userId, articleId, patch)

	return err
}
