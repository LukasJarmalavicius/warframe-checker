package cache

import (
	"log"
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

func (c *Cache) Unvaulted() []models.TrimmedItem {
	start := time.Now()
	c.mu.RLock()
	defer c.mu.RUnlock()
	items := make([]models.TrimmedItem, 0, len(c.data))
	for _, item := range c.data {
		if item.Vaulted || !item.IsPrime {
			continue
		}
		items = append(items, models.TrimmedItem{
			UniqueName: item.UniqueName,
			Name:       item.Name,
			Category:   item.Category,
			Vaulted:    item.Vaulted,
		})
	}
	log.Printf("cache: Unvaulted took %s", time.Since(start))
	return items
}
