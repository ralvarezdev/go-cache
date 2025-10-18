package go_cache

type (
	// Cache represents an in-memory cache
	Cache interface {
		Set(key string, value interface{}) error
		UpdateValue(key string, value interface{}) error
		Has(key string) bool
		Get(key string) (interface{}, bool)
		Delete(key string)
	}
)
