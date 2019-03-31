package cache

import (
	"github.com/elastic/beats/libbeat/common"
)

// TODO: Implement expiration
// TODO: Lock
type Cache struct {
	metadata map[string]common.MapStr
}

func NewCache() *Cache {
	return &Cache{
		metadata: make(map[string]common.MapStr),
	}
}

func (c *Cache) Get(key string) common.MapStr {
	return c.metadata[key]
}

func (c *Cache) Delete(key string) {
	delete(c.metadata, key)
}

func (c *Cache) Set(key string, data common.MapStr) {
	c.metadata[key] = data
}
