package abstractrepositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
)

type DigestsRepository interface {
	FindLatestDigestForToday(ctx context.Context, userId entities.UserId) (bool, string, error)
	AddDigest(ctx context.Context, userId entities.UserId, digest string) error
}
