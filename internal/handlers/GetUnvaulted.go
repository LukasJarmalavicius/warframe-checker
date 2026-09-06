package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"warframe-checker/internal/models"
)

type UnvaultedResponse struct {
	Warframes []models.WFCDItem
	Melee     []models.WFCDItem
	Primary   []models.WFCDItem
	Secondary []models.WFCDItem
}

type result struct {
	category string
	items    []models.WFCDItem
}

func fetchItems(link string) []models.WFCDItem {
	fetchTimer := time.Now()
	log.Printf("fetching %s\n", link)
	var items []models.WFCDItem
	if err := fetchJson(link, &items); err != nil {
		log.Println(err)
		return nil
	}

	var result []models.WFCDItem
	for _, item := range items {
		if item.Vaulted || !item.IsPrime {
			continue
		}
		result = append(result, item)
	}

	log.Printf("fetching took %v\n", time.Since(fetchTimer))
	return result
}

func (h *Handler) GetUnvaulted(w http.ResponseWriter, r *http.Request) {
	log.Println("/unvaulted")

	links := map[string]string{
		"Warframes": h.WFCD_JSON + "Warframes.json",
		"Primary":   h.WFCD_JSON + "Primary.json",
		"Secondary": h.WFCD_JSON + "Secondary.json",
		"Melee":     h.WFCD_JSON + "Melee.json",
	}

	results := make(chan result, len(links))
	var wg sync.WaitGroup
	for category, link := range links {
		wg.Add(1)
		go func(category, link string) {
			defer wg.Done()
			results <- result{category: category, items: fetchItems(link)}
		}(category, link)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var unvaulted UnvaultedResponse
	for r := range results {
		switch r.category {
		case "Warframes":
			unvaulted.Warframes = r.items
		case "Primary":
			unvaulted.Primary = r.items
		case "Secondary":
			unvaulted.Secondary = r.items
		case "Melee":
			unvaulted.Melee = r.items
		}
	}

	log.Printf("/unvaulted: returned %d items\n", len(unvaulted.Warframes)+len(unvaulted.Melee)+len(unvaulted.Primary)+len(unvaulted.Secondary))

	if err := json.NewEncoder(w).Encode(unvaulted); err != nil {
		log.Println(err)
		return
	}

}
