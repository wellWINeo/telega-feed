package abstractservices

import "context"

type PurgeService interface {
	Purge(ctx context.Context) error
}
