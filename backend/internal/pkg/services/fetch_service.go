package services

import (
	"TelegaFeed/internal/pkg/core/abstractions/infrastructure/providers"
	"TelegaFeed/internal/pkg/core/entities"
	"TelegaFeed/pkg/myhttp"
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

type FetchService struct {
	providersMap map[entities.FeedType]abstractproviders.FetchProvider
	http         myhttp.HttpClient
}

func NewFetchService(providersMap map[string]abstractproviders.FetchProvider) *FetchService {
	return &FetchService{
		providersMap: providersMap,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func NewFetchServiceWithHttpClient(
	providersMap map[string]abstractproviders.FetchProvider,
	http myhttp.HttpClient,
) *FetchService {
	return &FetchService{
		providersMap: providersMap,
		http:         http,
	}
}

func (f FetchService) DetectType(feedSource *entities.FeedSource) entities.FeedType {
	resp, err := myhttp.Fetch(f.http, http.MethodGet, feedSource.FeedUrl)
	if err != nil {
		return entities.DefaultFeedType
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	for k, v := range f.providersMap {
		right, err := v.CheckType(resp.Body)
		if err != nil {
			continue
		}

		if right {
			return k
		}
	}

	return entities.DefaultFeedType
}

func (f FetchService) FetchArticles(ctx context.Context, feedSources []*entities.FeedSource) []*entities.Article {
	wg := sync.WaitGroup{}
	articlesChan := make(chan *entities.Article)

	wg.Add(len(feedSources))

	for i := range feedSources {
		go f.fetchArticlesBySource(ctx, feedSources[i], &wg, articlesChan)
	}

	go func() {
		wg.Wait()
		close(articlesChan)
	}()

	articles := make([]*entities.Article, 0)
	for article := range articlesChan {
		articles = append(articles, article)
	}

	return articles
}

func (f FetchService) fetchArticlesBySource(ctx context.Context, source *entities.FeedSource, wg *sync.WaitGroup, out chan<- *entities.Article) {
	defer wg.Done()

	provider := f.providersMap[entities.DefaultFeedType]

	articles, err := provider.FetchArticles(ctx, source)
	if err != nil {
		log.Printf("Failed to fetch articles from %s, error: %v", source.FeedUrl, err)
	}

	for _, article := range articles {
		out <- article
	}
}
