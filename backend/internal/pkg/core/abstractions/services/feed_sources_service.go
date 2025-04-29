package abstractservices

import abstractrepositories "TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"

// FeedSourcesService just forwards interface repository
type FeedSourcesService interface {
	abstractrepositories.FeedSourceRepository
}
