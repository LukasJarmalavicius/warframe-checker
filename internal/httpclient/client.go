package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	limiter    *rate.Limiter
	httpClient *http.Client
}

func NewClient(limiter *rate.Limiter, httpClient *http.Client) *Client {
	return &Client{
		limiter:    limiter,
		httpClient: httpClient,
	}
}

func NewDefaultClient() *Client {
	return NewClient(rate.NewLimiter(rate.Limit(3), 1), &http.Client{
		Timeout: 10 * time.Second,
	})
}

func FetchJson(url string, target any) error {
	client := NewDefaultClient()
	return client.FetchJson(url, target)
}

func (c *Client) FetchJson(url string, target any) error {
	if err := c.limiter.Wait(context.Background()); err != nil {
		return fmt.Errorf("fetching %s: %w", url, err)
	}

	start := time.Now()
	log.Printf("fetching %s\n", url)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("fetching %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetching %s: unexpected status code: %d", url, resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decoding %s: %w", url, err)
	}
	log.Printf("fetching took %dms\n", time.Since(start).Milliseconds())
	return nil
}
