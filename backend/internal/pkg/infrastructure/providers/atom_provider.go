package providers

import (
	"TelegaFeed/internal/pkg/core/entities"
	"TelegaFeed/pkg/myhttp"
	"TelegaFeed/pkg/mytime"
	"TelegaFeed/pkg/parser/atom"
	"TelegaFeed/pkg/parser/feedutil"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const atomSpecUrl = "www.w3.org/2005/Atom"

type AtomProvider struct {
	http myhttp.HttpClient
}

func NewAtomProvider() *AtomProvider {
	return &AtomProvider{
		http: &http.Client{Timeout: 15 * time.Second},
	}
}

func NewAtomProviderWithClient(http myhttp.HttpClient) *AtomProvider {
	return &AtomProvider{http: http}
}

func (a AtomProvider) CheckType(reader io.Reader) (bool, error) {
	startElement, err := getXMLStartElement(reader)
	if err != nil {
		return false, err
	}

	if startElement.Name.Local != "feed" {
		return false, nil
	}

	for _, attr := range startElement.Attr {
		if attr.Name.Local == "xmlns" && strings.Contains(attr.Value, atomSpecUrl) {
			return true, nil
		}
	}

	return false, nil
}

func (a AtomProvider) FetchArticles(ctx context.Context, feedSource *entities.FeedSource) ([]*entities.Article, error) {
	resp, err := myhttp.Fetch(a.http, http.MethodGet, feedSource.FeedUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch articles from %s: %w", feedSource.FeedUrl, err)
	}

	defer resp.Body.Close()

	atomFeed, err := atom.ParseAtom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	articles := make([]*entities.Article, len(atomFeed.Entries))
	for i, entry := range atomFeed.Entries {
		articles[i] = &entities.Article{
			Title:       strings.TrimSpace(entry.Title),
			Text:        strings.TrimSpace(entry.Content.Value),
			Url:         atom.GetLink(entry.Links),
			PreviewUrl:  atom.GetPreviewLink(entry.Links),
			AddedAt:     mytime.NowUTC(),
			PublishedAt: feedutil.ParseFeedDate(entry.Updated),
		}
	}

	return articles, nil
}
