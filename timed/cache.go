package timed

import (
	gocache "github.com/ralvarezdev/go-cache"
	"sync"
	"time"
)

type (
	// Item represents a cached value with an expiration time
	Item struct {
		value     interface{}
		expiresAt time.Time
	}

	// Cache represents an in-memory cache
	Cache struct {
		items map[string]*Item
		mutex sync.RWMutex
	}
)

// NewItem creates a new cache item
func NewItem(value interface{}, expiresAt time.Time) *Item {
	return &Item{
		expiresAt: expiresAt,
		value:     value,
	}
}

// Value returns the cached value
func (i *Item) Value() interface{} {
	return i.value
}

// ExpiresAt returns the expiration time of the cached value
func (i *Item) ExpiresAt() time.Time {
	return i.expiresAt
}

// Expired returns true if the cached value has expired
func (i *Item) Expired() bool {
	return time.Now().After(i.expiresAt)
}

// NewCache creates a new Cache instance
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]*Item),
	}
}

// Set adds the item to the cache
func (c *Cache) Set(key string, item *Item) error {
	// Lock the cache
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if the item is nil or has expired
	if item == nil {
		return gocache.ErrNilItem
	}
	if item.Expired() {
		return ErrExpiredItem
	}

	// Add the item to the cache
	c.items[key] = item
	return nil
}

// Has checks if the cache contains a key
func (c *Cache) Has(key string) bool {
	// Lock the cache
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Check if the item exists
	item, found := c.items[key]
	if !found {
		return false
	}

	// Check if the item has expired, and remove it if it has
	if item.Expired() {
		delete(c.items, key)
		return false
	}
	return true
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	// Lock the cache
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Check if the item exists
	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// Check if the item has expired, and remove it if it has
	if item.Expired() {
		delete(c.items, key)
		return nil, false
	}

	return item.value, true
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) {
	// Lock the cache
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
}
