package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	API_URL   string
	WFCD_JSON string
	WFCD_API  string
}

func fetchJson(url string, target any) error {
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
	return nil
}
