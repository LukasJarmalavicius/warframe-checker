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

var limiter = rate.NewLimiter(rate.Limit(3), 1)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func FetchJson(url string, target any) error {
	if err := limiter.Wait(context.Background()); err != nil {
		return fmt.Errorf("fetching %s: %w", url, err)
	}

	start := time.Now()
	log.Printf("fetching %s\n", url)
	resp, err := httpClient.Get(url)
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
