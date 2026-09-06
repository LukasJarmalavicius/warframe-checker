package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"warframe-checker/internal/models"
)

type UnvaultedResponse struct {
	Warframes []models.WFCDItem
	Melee     []models.WFCDItem
	Primary   []models.WFCDItem
	Secondary []models.WFCDItem
}

func (h *Handler) GetUnvaulted(w http.ResponseWriter, r *http.Request) {
	log.Println("/unvaulted")

	links := []string{h.WFCD_JSON + "Warframes.json", h.WFCD_JSON + "Melee.json", h.WFCD_JSON + "Primary.json", h.WFCD_JSON + "Secondary.json"}
	var unvaulted UnvaultedResponse

	i := 1
	for _, link := range links {
		log.Printf("fetching %s\n", link)
		resp, err := http.Get(link)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		var items []models.WFCDItem
		if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
			log.Println(err)
			return
		}

		for _, item := range items {
			if item.Vaulted || !item.IsPrime {
				continue
			}
			switch i {
			case 1:
				unvaulted.Warframes = append(unvaulted.Warframes, item)
			case 2:
				unvaulted.Melee = append(unvaulted.Melee, item)
			case 3:
				unvaulted.Primary = append(unvaulted.Primary, item)
			case 4:
				unvaulted.Secondary = append(unvaulted.Secondary, item)
			}
		}
		i++
	}

	log.Printf("/unvaulted: returned %d items\n", len(unvaulted.Warframes)+len(unvaulted.Melee)+len(unvaulted.Primary)+len(unvaulted.Secondary))

	if err := json.NewEncoder(w).Encode(unvaulted); err != nil {
		log.Println(err)
		return
	}

}
