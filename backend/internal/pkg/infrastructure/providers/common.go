package providers

import (
	"TelegaFeed/pkg/myhttp"
	"encoding/xml"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func findPreviewUrl(articleUrl string, httpClient myhttp.HttpClient) string {
	parsedUrl, err := url.Parse(articleUrl)
	if err != nil {
		return ""
	}

	resp, err := httpClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    parsedUrl,
	})
	if err != nil || resp.StatusCode != http.StatusOK {
		return ""
	}

	defer resp.Body.Close()

	utf8Reader, err := charset.NewReader(resp.Body, "text/html")
	if err != nil {
		utf8Reader = resp.Body
	}

	ogImageUrl, err := findImageUrlFromOpenGraph(utf8Reader)
	if err != nil {
		// fallback to favicon
		return fmt.Sprintf("%s://%s/favicon.ico", parsedUrl.Scheme, parsedUrl.Host)
	}

	return ogImageUrl
}

func findImageUrlFromOpenGraph(body io.Reader) (string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return "", err
	}

	var walkFunc func(n *html.Node) string
	walkFunc = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "meta" {
			var property, content string
			for _, attr := range n.Attr {
				switch attr.Key {
				case "property":
					property = attr.Val
				case "content":
					content = attr.Val
				}
			}

			if strings.HasPrefix(property, "og:image") && content != "" {
				return content
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if result := walkFunc(c); result != "" {
				return result
			}
		}

		return ""
	}

	if ogImageUrl := walkFunc(doc); ogImageUrl != "" {
		return ogImageUrl, nil
	}

	return "", errors.New("could not find og image")
}

func checkXMLStartElement[T any](reader io.Reader, checkFunc func(el xml.StartElement) T) (T, error) {
	decoder := xml.NewDecoder(reader)

	for {
		token, err := decoder.Token()
		if err != nil {
			return *new(T), err
		}

		if startElement, ok := token.(xml.StartElement); ok {
			return checkFunc(startElement), nil
		}
	}
}

func getXMLStartElement(reader io.Reader) (*xml.StartElement, error) {
	decoder := xml.NewDecoder(reader)

	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}

		if startElement, ok := token.(xml.StartElement); ok {
			return &startElement, nil
		}
	}
}
