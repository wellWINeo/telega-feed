package providers

import (
	"TelegaFeed/internal/pkg/core/entities"
	"TelegaFeed/pkg/myhttp"
	"TelegaFeed/pkg/mytime"
	"TelegaFeed/pkg/parser/feedutil"
	"TelegaFeed/pkg/parser/rdf"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const rdfSpecUrl = "www.w3.org/1999/02/22-rdf-syntax-ns"

type RDFProvider struct {
	http myhttp.HttpClient
}

func NewRDFProvider() *RDFProvider {
	return &RDFProvider{
		http: &http.Client{Timeout: 15 * time.Second},
	}
}

func NewRDFProviderWithClient(http myhttp.HttpClient) *RDFProvider {
	return &RDFProvider{http: http}
}

func (r *RDFProvider) CheckType(reader io.Reader) (bool, error) {
	startElement, err := getXMLStartElement(reader)
	if err != nil {
		return false, err
	}

	for _, attr := range startElement.Attr {
		if attr.Name.Local == "xmlns" && strings.Contains(attr.Value, rdfSpecUrl) {
			return true, nil
		}
	}

	return false, nil
}

func (r *RDFProvider) FetchArticles(ctx context.Context, feedSource *entities.FeedSource) ([]*entities.Article, error) {
	resp, err := myhttp.Fetch(r.http, http.MethodGet, feedSource.FeedUrl)
	if err != nil {
		return nil, fmt.Errorf("fetch articles failed: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	rdfFeed, err := rdf.ParseRDF(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	articles := make([]*entities.Article, len(rdfFeed.Items))
	for i, item := range rdfFeed.Items {
		articles[i] = &entities.Article{
			Title:       strings.TrimSpace(item.Title),
			Text:        strings.TrimSpace(item.Description),
			Url:         item.Link,
			PreviewUrl:  findPreviewUrl(item.Link, r.http),
			AddedAt:     mytime.NowUTC(),
			PublishedAt: feedutil.ParseFeedDate(item.Date),
		}
	}

	return articles, nil
}
