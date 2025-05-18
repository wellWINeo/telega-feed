package abstractrepositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
	"time"
)

type DigestsRepository interface {
	FindLatestDigestForToday(ctx context.Context, userId entities.UserId) (bool, string, error)
	AddDigest(ctx context.Context, userId entities.UserId, digest string) error
	DeleteDigestsGeneratedBefore(ctx context.Context, datetime time.Time) error
}
