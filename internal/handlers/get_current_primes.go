package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type UnvaultedResponse struct {
	Warframes []string
	Primary   []string
	Secondary []string
	Melee     []string
}

func (h *Handler) GetCurrentPrimes(w http.ResponseWriter, r *http.Request) {
	log.Println("/currentPrimes")

	start := time.Now()
	items := h.inventoryService.GetCurrentPrimes()

	var unvaulted UnvaultedResponse
	for _, item := range items {
		switch item.Category {
		case "Warframes":
			unvaulted.Warframes = append(unvaulted.Warframes, item.Name)
		case "Melee":
			unvaulted.Melee = append(unvaulted.Melee, item.Name)
		case "Primary":
			unvaulted.Primary = append(unvaulted.Primary, item.Name)
		case "Secondary":
			unvaulted.Secondary = append(unvaulted.Secondary, item.Name)
		}
	}

	log.Printf("/currentPrimes: returned %d items\n", len(items))
	log.Printf("cache: GetCurrentPrimes took %s", time.Since(start))

	if err := json.NewEncoder(w).Encode(unvaulted); err != nil {
		log.Println(err)
		return
	}

}
