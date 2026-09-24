package sep

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
}

func NewFetcher() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (c *Client) Fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set(
		"User-Agent",
		"sep2epub/0.1 personal EPUB converter",
	)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch SEP: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("unexpected HTTP: %s", res.Status)
	}

	return res.Body, nil
}
