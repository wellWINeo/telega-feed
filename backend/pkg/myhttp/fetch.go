package myhttp

import (
	"fmt"
	"net/http"
	"net/url"
)

func Fetch(client HttpClient, method string, requestUrl string) (*http.Response, error) {
	parsedUrl, err := url.Parse(requestUrl)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(&http.Request{
		Method: method,
		URL:    parsedUrl,
	})

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code does not indicate success: %s", resp.Status)
	}

	return resp, nil
}
