package timed

import (
	gocache "github.com/ralvarezdev/go-cache"
	"sync"
	"time"
)

type (
	// Item represents a cached value with an expiration time
	Item struct {
		expiresAt time.Time
		value     interface{}
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
		expiresAt,
		value,
	}
}

// GetValue returns the cached value
func (i *Item) GetValue() interface{} {
	return i.value
}

// SetValue sets the cached value
func (i *Item) SetValue(value interface{}) {
	i.value = value
}

// GetExpiresAt returns the expiration time of the cached value
func (i *Item) GetExpiresAt() time.Time {
	return i.expiresAt
}

// SetExpiresAt sets the expiration time of the cached value
func (i *Item) SetExpiresAt(expiresAt time.Time) {
	i.expiresAt = expiresAt
}

// HasExpired returns true if the cached value has expired
func (i *Item) HasExpired() bool {
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
	if item.HasExpired() {
		return ErrItemHasExpired
	}

	// Add the item to the cache
	c.items[key] = item
	return nil
}

// UpdateValue updates the value of an item in the cache
func (c *Cache) UpdateValue(key string, value interface{}) error {
	// Lock the cache
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if the item exists
	item, found := c.items[key]
	if !found {
		return gocache.ErrItemNotFound
	}

	// Update the value
	item.value = value
	return nil
}

// UpdateExpiresAt updates the expiration time of an item in the cache
func (c *Cache) UpdateExpiresAt(key string, expiresAt time.Time) error {
	// Lock the cache
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if the item exists
	item, found := c.items[key]
	if !found {
		return gocache.ErrItemNotFound
	}

	// Update the expiration time
	item.expiresAt = expiresAt
	return nil
}

// Has checks if the cache contains a key
func (c *Cache) Has(key string) bool {
	// Get the item from the cache
	_, found := c.Get(key)
	return found
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
	if item.HasExpired() {
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
