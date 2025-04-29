package providers

import (
	"TelegaFeed/internal/pkg/core/entities"
	"TelegaFeed/pkg/myhttp"
	"TelegaFeed/pkg/mytime"
	"TelegaFeed/pkg/parser/feedutil"
	"TelegaFeed/pkg/parser/rss"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type RssProvider struct {
	http myhttp.HttpClient
}

func NewRssProvider() *RssProvider {
	return &RssProvider{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func NewRssProviderWithClient(http myhttp.HttpClient) *RssProvider {
	return &RssProvider{
		http: http,
	}
}

func (r *RssProvider) CheckType(reader io.Reader) (bool, error) {
	startElement, err := getXMLStartElement(reader)
	if err != nil {
		return false, err
	}

	return startElement.Name.Local == "rss", nil
}

func (r *RssProvider) FetchArticles(ctx context.Context, feedSource *entities.FeedSource) ([]*entities.Article, error) {
	resp, err := myhttp.Fetch(r.http, http.MethodGet, feedSource.FeedUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch articles from %s, error: %v", feedSource.FeedUrl, err)
	}

	defer resp.Body.Close()

	rssFeed, err := rss.ParseRSS(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rss: %w", err)
	}

	articles := make([]*entities.Article, len(rssFeed.Channel.Items))
	for i, channelItem := range rssFeed.Channel.Items {
		articles[i] = &entities.Article{
			Title:       strings.TrimSpace(channelItem.Title),
			Text:        strings.TrimSpace(channelItem.Description),
			Url:         channelItem.Link,
			PreviewUrl:  rss.GetPreviewUrlFromRSS(channelItem.Enclosures),
			AddedAt:     mytime.NowUTC(),
			PublishedAt: feedutil.ParseFeedDate(channelItem.PubDate),
		}
	}

	return articles, nil
}
