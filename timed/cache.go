package timed

import (
	"sync"
	"time"

	gocache "github.com/ralvarezdev/go-cache"
)

type (
	// TimedItem represents a cached value with an expiration time
	TimedItem struct {
		expiresAt time.Time
		value     interface{}
	}

	// DefaultTimedCache represents an in-memory cache
	DefaultTimedCache struct {
		items map[string]*TimedItem
		mutex sync.RWMutex
	}
)

// NewTimedItem creates a new cache item
//
// Parameters:
//
//   - value: The value to be cached
//   - expiresAt: The expiration time of the cached value
//
// Returns:
//
//   - *TimedItem: A pointer to the newly created cache item
func NewTimedItem(value interface{}, expiresAt time.Time) *TimedItem {
	return &TimedItem{
		expiresAt,
		value,
	}
}

// GetValue returns the cached value
//
// Returns:
//
//   - interface{}: The cached value
func (i *TimedItem) GetValue() interface{} {
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
func (i *TimedItem) SetValue(value interface{}) {
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
func (i *TimedItem) GetExpiresAt() time.Time {
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
func (i *TimedItem) SetExpiresAt(expiresAt time.Time) {
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
func (i *TimedItem) HasExpired() bool {
	if i == nil {
		return true
	}
	return time.Now().After(i.expiresAt)
}

// NewDefaultTimedCache creates a new DefaultTimedCache instance
func NewDefaultTimedCache() *DefaultTimedCache {
	return &DefaultTimedCache{
		items: make(map[string]*TimedItem),
	}
}

// Set adds the item to the cache
//
// Parameters:
//
//   - key: The key to associate with the cached value
//   - value: The item to be cached
//
// Returns:
//
//   - error: An error if the item is nil or has expired
func (d *DefaultTimedCache) Set(key string, value interface{}) error {
	if d == nil {
		return gocache.ErrNilCache
	}

	// Lock the cache
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Type assert the value to *TimedItem
	item, ok := value.(*TimedItem)
	if !ok {
		return ErrValueMustBeATimedItem
	}

	// Check if the item is nil or has expired
	if item == nil {
		return gocache.ErrNilItem
	}
	if item.HasExpired() {
		return ErrItemHasExpired
	}

	// Add the item to the cache
	d.items[key] = item
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
func (d *DefaultTimedCache) UpdateValue(key string, value interface{}) error {
	if d == nil {
		return gocache.ErrNilCache
	}

	// Lock the cache
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Check if the item exists
	item, found := d.items[key]
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
func (d *DefaultTimedCache) UpdateExpiresAt(
	key string,
	expiresAt time.Time,
) error {
	if d == nil {
		return gocache.ErrNilCache
	}

	// Lock the cache
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Check if the item exists
	item, found := d.items[key]
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
func (d *DefaultTimedCache) Has(key string) bool {
	if d == nil {
		return false
	}

	// Get the item from the cache
	_, found := d.Get(key)
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
func (d *DefaultTimedCache) Get(key string) (interface{}, bool) {
	if d == nil {
		return nil, false
	}

	// Lock the cache
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	// Check if the item exists
	item, found := d.items[key]
	if !found {
		return nil, false
	}

	// Check if the item has expired, and remove it if it has
	if item.HasExpired() {
		delete(d.items, key)
		return nil, false
	}

	return item.value, true
}

// Delete removes a value from the cache
//
// Parameters:
//
//   - key: The key to remove from the cache
func (d *DefaultTimedCache) Delete(key string) {
	if d == nil {
		return
	}

	// Lock the cache
	d.mutex.Lock()
	defer d.mutex.Unlock()

	delete(d.items, key)
}

// GetExpirationTime retrieves the expiration time of a cached item
//
// Parameters:
//
//   - key: The key associated with the cached value
//
// Returns:
//
//   - time.Time: The expiration time of the cached item, or zero time if not found
func (d *DefaultTimedCache) GetExpirationTime(key string) time.Time {
	if d == nil {
		return time.Time{}
	}

	// Lock the cache
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	// Check if the item exists
	item, found := d.items[key]
	if !found {
		return time.Time{}
	}

	return item.expiresAt
}

// UpdateExpirationTime updates the expiration time of a cached item
//
// Parameters:
//
//   - key: The key associated with the cached value
//   - expiresAt: The new expiration time to be set
//
// Returns:
//
//   - error: An error if the item is not found
func (d *DefaultTimedCache) UpdateExpirationTime(
	key string,
	expiresAt time.Time,
) error {
	if d == nil {
		return gocache.ErrNilCache
	}

	// Lock the cache
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Check if the item exists
	item, found := d.items[key]
	if !found {
		return gocache.ErrItemNotFound
	}

	// Update the expiration time
	item.expiresAt = expiresAt
	return nil
}
