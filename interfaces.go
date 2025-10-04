package go_cache

type (
	// Item represents a cached value
	Item interface {
		GetValue() interface{}
		SetValue(value interface{})
	}

	// Cache represents an in-memory cache
	Cache interface {
		Set(key string, item Item) error
		UpdateValue(key string, value interface{}) error
		Has(key string) bool
		Get(key string) (Item, bool)
		Delete(key string)
	}
)
