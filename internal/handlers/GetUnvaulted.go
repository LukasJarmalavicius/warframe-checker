package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"warframe-checker/internal/models"
)

type UnvaultedResponse struct {
	Warframes []models.TrimmedItem
	Melee     []models.TrimmedItem
	Primary   []models.TrimmedItem
	Secondary []models.TrimmedItem
}

func (h *Handler) GetUnvaulted(w http.ResponseWriter, r *http.Request) {
	log.Println("/unvaulted")

	start := time.Now()
	items := h.cache.Unvaulted()

	var unvaulted UnvaultedResponse
	for _, item := range items {
		switch item.Category {
		case "Warframes":
			unvaulted.Warframes = append(unvaulted.Warframes, models.TrimmedItem(item))
		case "Melee":
			unvaulted.Melee = append(unvaulted.Melee, models.TrimmedItem(item))
		case "Primary":
			unvaulted.Primary = append(unvaulted.Primary, models.TrimmedItem(item))
		case "Secondary":
			unvaulted.Secondary = append(unvaulted.Secondary, models.TrimmedItem(item))
		}
	}

	log.Printf("/unvaulted: returned %d items\n", len(unvaulted.Warframes)+len(unvaulted.Melee)+len(unvaulted.Primary)+len(unvaulted.Secondary))
	log.Printf("cache: GetUnvaulted took %s", time.Since(start))

	if err := json.NewEncoder(w).Encode(unvaulted); err != nil {
		log.Println(err)
		return
	}

}
