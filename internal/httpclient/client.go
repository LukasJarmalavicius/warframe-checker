package httpclient

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func FetchJson(url string, target any) error {
	start := time.Now()
	log.Printf("fetching %s\n", url)
	resp, err := http.Get(url)
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
