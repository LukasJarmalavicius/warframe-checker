package cache

import (
	"log"
	"sort"
	"strings"
	"sync"
	"time"
	"warframe-checker/internal/httpclient"
	"warframe-checker/internal/models"
)

type Cache struct {
	data map[string]models.WFCDItem
	mu   sync.RWMutex
}

type itemResult struct {
	items []models.WFCDItem
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]models.WFCDItem),
	}
}
func (c *Cache) Set(key string, value models.WFCDItem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}
func (c *Cache) Get(key string) (models.WFCDItem, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

func (c *Cache) Load(wfcdJSONbase string) error {
	links := []string{
		wfcdJSONbase + "Warframes.json",
		wfcdJSONbase + "Melee.json",
		wfcdJSONbase + "Primary.json",
		wfcdJSONbase + "Secondary.json",
	}

	results := make(chan itemResult, len(links))
	var wg sync.WaitGroup

	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			var items []models.WFCDItem
			if err := httpclient.FetchJson(link, &items); err != nil {
				log.Println(err)
				return
			}
			results <- itemResult{items: items}
		}(link)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	newData := make(map[string]models.WFCDItem)
	for r := range results {
		for _, item := range r.items {
			newData[strings.ToLower(item.Name)] = item
		}
	}

	c.mu.Lock()
	c.data = newData
	c.mu.Unlock()
	log.Printf("cache: loaded %d items", len(newData))
	return nil
}

func (c *Cache) All() []models.WFCDItem {
	start := time.Now()
	c.mu.RLock()
	defer c.mu.RUnlock()
	items := make([]models.WFCDItem, 0, len(c.data))
	for _, item := range c.data {
		items = append(items, item)
	}
	log.Printf("cache: All took %s", time.Since(start))
	return items
}

type PartialSet struct {
	SetName      string   `json:"setName"`
	Owned        []string `json:"owned"`
	Missing      []string `json:"missing"`
	TotalParts   int      `json:"totalParts"`
	MissingCount int      `json:"missingCount"`
}

func (c *Cache) AlmostCompleteSets(owned []string, maxMissing int) []PartialSet {
	start := time.Now()
	ownedSet := make(map[string]bool, len(owned))
	for _, item := range owned {
		ownedSet[strings.ToLower(strings.TrimSpace(item))] = true
	}

	var result []PartialSet
	for _, item := range c.All() {
		if !item.IsPrime || len(item.Components) == 0 {
			continue
		}

		var have, missing []string

		for _, component := range item.Components {
			if component.Type == "Resource" {
				continue
			}
			fullName := item.Name + " " + component.Name
			if ownedSet[strings.ToLower(fullName)] {
				have = append(have, fullName)
			} else {
				missing = append(missing, fullName)
			}
		}

		if len(have) > 0 && len(missing) > 0 && len(missing) <= maxMissing {
			result = append(result, PartialSet{
				SetName:      item.Name,
				Owned:        have,
				Missing:      missing,
				TotalParts:   len(item.Components),
				MissingCount: len(missing),
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].MissingCount < result[j].MissingCount
	})

	log.Printf("cache: AlmostCompleteSets took %s", time.Since(start))
	return result
}
