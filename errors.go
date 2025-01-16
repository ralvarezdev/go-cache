package go_cache

import (
	"errors"
)

var (
	ErrNilCache = errors.New("cache cannot be nil")
	ErrNilItem  = errors.New("item cannot be nil")
)
