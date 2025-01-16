package go_cache

type (
	// Item represents a cached value
	Item interface {
		Value() interface{}
	}

	// Cache represents an in-memory cache
	Cache interface {
		Set(key string, item Item) error
		Has(key string) bool
		Get(key string) (Item, bool)
		Delete(key string)
	}
)
