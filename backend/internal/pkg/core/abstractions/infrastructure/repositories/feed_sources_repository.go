package abstractrepositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
)

type FeedSourceRepository interface {
	AddSource(ctx context.Context, userId entities.UserId, source *entities.FeedSource) error
	GetSources(ctx context.Context, userId entities.UserId) ([]*entities.FeedSource, error)
	GetSource(ctx context.Context, userId entities.UserId, sourceId entities.FeedSourceId) (*entities.FeedSource, error)
	UpdateSource(ctx context.Context, userId entities.UserId, sourceId entities.FeedSourceId, patch *entities.FeedSourcePatch) error
	DeleteSource(ctx context.Context, userId entities.UserId, sourceId entities.FeedSourceId) error
}
