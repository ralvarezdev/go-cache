package timed

import (
	"sync"
	"time"

	gocache "github.com/ralvarezdev/go-cache"
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
//
// Parameters:
//
//   - value: The value to be cached
//   - expiresAt: The expiration time of the cached value
//
// Returns:
//
//   - *Item: A pointer to the newly created cache item
func NewItem(value interface{}, expiresAt time.Time) *Item {
	return &Item{
		expiresAt,
		value,
	}
}

// GetValue returns the cached value
//
// Returns:
//
//   - interface{}: The cached value
func (i *Item) GetValue() interface{} {
	if i == nil {
		return nil
	}
	return i.value
}

// SetValue sets the cached value
//
// Parameters:
//
//   - value: The value to be cached
func (i *Item) SetValue(value interface{}) {
	if i == nil {
		return
	}
	i.value = value
}

// GetExpiresAt returns the expiration time of the cached value
//
// Returns:
//
//   - time.Time: The expiration time of the cached value
func (i *Item) GetExpiresAt() time.Time {
	if i == nil {
		return time.Time{}
	}
	return i.expiresAt
}

// SetExpiresAt sets the expiration time of the cached value
//
// Parameters:
//
//   - expiresAt: The expiration time to be set
func (i *Item) SetExpiresAt(expiresAt time.Time) {
	if i == nil {
		return
	}
	i.expiresAt = expiresAt
}

// HasExpired returns true if the cached value has expired
//
// Returns:
//
//   - bool: True if the cached value has expired, false otherwise
func (i *Item) HasExpired() bool {
	if i == nil {
		return true
	}
	return time.Now().After(i.expiresAt)
}

// NewCache creates a new Cache instance
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]*Item),
	}
}

// Set adds the item to the cache
//
// Parameters:
//
//   - key: The key to associate with the cached value
//   - item: The item to be cached
//
// Returns:
//
//   - error: An error if the item is nil or has expired
func (c *Cache) Set(key string, item *Item) error {
	if c == nil {
		return gocache.ErrNilCache
	}

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
//
// Parameters:
//
//   - key: The key associated with the cached value
//   - value: The new value to be set
//
// Returns:
//
//   - error: An error if the item is not found
func (c *Cache) UpdateValue(key string, value interface{}) error {
	if c == nil {
		return gocache.ErrNilCache
	}

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
//
// Parameters:
//
//   - key: The key associated with the cached value
//   - expiresAt: The new expiration time to be set
//
// Returns:
//
//   - error: An error if the item is not found
func (c *Cache) UpdateExpiresAt(key string, expiresAt time.Time) error {
	if c == nil {
		return gocache.ErrNilCache
	}

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
//
// Parameters:
//
//   - key: The key to check in the cache
//
// Returns:
//
//   - bool: True if the key exists in the cache and has not expired, false otherwise
func (c *Cache) Has(key string) bool {
	if c == nil {
		return false
	}

	// Get the item from the cache
	_, found := c.Get(key)
	return found
}

// Get retrieves a value from the cache
//
// Parameters:
//
//   - key: The key to retrieve from the cache
//
// Returns:
//
//   - interface{}: The cached value, or nil if not found or expired
//   - bool: True if the value was found and not expired, false otherwise
func (c *Cache) Get(key string) (interface{}, bool) {
	if c == nil {
		return nil, false
	}

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
//
// Parameters:
//
//   - key: The key to remove from the cache
func (c *Cache) Delete(key string) {
	if c == nil {
		return
	}

	// Lock the cache
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
}
