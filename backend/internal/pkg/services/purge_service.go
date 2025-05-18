package services

import (
	abstractrepositories "TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"
	"context"
	"log"
	"time"
)

const DeletionDays = 30

type PurgeService struct {
	feedSourceRepository abstractrepositories.FeedSourceRepository
	feedRepository       abstractrepositories.FeedRepository
	summariesRepository  abstractrepositories.SummariesRepository
	digestsRepository    abstractrepositories.DigestsRepository
}

func (p *PurgeService) Purge(ctx context.Context) error {
	deletionDate := time.Now().UTC().AddDate(0, 0, -DeletionDays)

	// purge feed sources
	if err := p.feedSourceRepository.DeleteOrphanedSources(ctx); err != nil {
		log.Printf("failed to delete orphaned sources: %v", err)
	}

	// purge articles
	if err := p.feedRepository.DeleteOrphanedArticles(ctx); err != nil {
		log.Printf("failed to delete orphaned articles: %v", err)
	}

	if err := p.feedRepository.DeleteArticlesAddedBefore(ctx, deletionDate); err != nil {
		log.Printf("failed to delete old articles: %v", err)
	}

	// purge summaries
	if err := p.summariesRepository.DeleteOrphanedSummaries(ctx); err != nil {
		log.Printf("failed to delete old summaries: %v", err)
	}

	// purge old digests
	if err := p.digestsRepository.DeleteDigestsGeneratedBefore(ctx, deletionDate); err != nil {
		log.Printf("failed to delete digests: %v", err)
	}

	return nil
}
