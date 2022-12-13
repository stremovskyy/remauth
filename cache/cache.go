package cache

import (
	"math"
	"sync"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, duration time.Duration)
	Get(key string) (interface{}, bool)
}

type cache struct {
	sync.Mutex
	elementsCount     int64
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]item
}

type item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

func New(defaultExpiration, cleanupInterval time.Duration) Cache {
	items := make(map[string]item)

	cache := cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.startGC()
	}

	return &cache
}

func (c *cache) Set(key string, value interface{}, duration time.Duration) {
	if c.elementsCount > math.MaxUint16 {
		return
	}

	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.items[key] = item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}

	c.elementsCount++
}

func (c *cache) Get(key string) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	if item, found := c.items[key]; found {
		if item.Expiration > 0 {
			if time.Now().UnixNano() > item.Expiration {
				return false, false
			}

			return item.Value, true
		}
	}

	return false, false
}

func (c *cache) startGC() {
	go c.gC()
}

func (c *cache) gC() {

	for {
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}

func (c *cache) expiredKeys() (keys []string) {
	c.Lock()
	defer c.Unlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

func (c *cache) clearItems(keys []string) {
	c.Lock()
	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
		c.elementsCount--
	}
}
