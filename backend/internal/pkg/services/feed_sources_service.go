package services

import (
	abstractrepositories "TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"TelegaFeed/internal/pkg/core/entities"
	"context"
)

type FeedSourcesService struct {
	fetchService         abstractservices.FetchService
	feedSourceRepository abstractrepositories.FeedSourceRepository
}

func NewFeedSourcesService(
	fetchService abstractservices.FetchService,
	feedSourceRepository abstractrepositories.FeedSourceRepository,
) *FeedSourcesService {
	return &FeedSourcesService{
		fetchService:         fetchService,
		feedSourceRepository: feedSourceRepository,
	}
}

func (f *FeedSourcesService) AddSource(ctx context.Context, userId entities.UserId, source *entities.FeedSource) error {
	source.Type = f.fetchService.DetectType(source)

	return f.feedSourceRepository.AddSource(ctx, userId, source)
}

func (f *FeedSourcesService) GetSources(ctx context.Context, userId entities.UserId) ([]*entities.FeedSource, error) {
	return f.feedSourceRepository.GetSources(ctx, userId)
}

func (f *FeedSourcesService) GetSource(
	ctx context.Context,
	userId entities.UserId,
	sourceId entities.FeedSourceId,
) (*entities.FeedSource, error) {
	return f.feedSourceRepository.GetSource(ctx, userId, sourceId)
}

func (f *FeedSourcesService) UpdateSource(
	ctx context.Context,
	userId entities.UserId,
	sourceId entities.FeedSourceId,
	patch *entities.FeedSourcePatch,
) error {
	return f.feedSourceRepository.UpdateSource(ctx, userId, sourceId, patch)
}

func (f *FeedSourcesService) DeleteSource(
	ctx context.Context,
	userId entities.UserId,
	sourceId entities.FeedSourceId,
) error {
	return f.feedSourceRepository.DeleteSource(ctx, userId, sourceId)
}
