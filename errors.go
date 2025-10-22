package gocache

import (
	"errors"
)

var (
	ErrNilCache     = errors.New("cache cannot be nil")
	ErrNilItem      = errors.New("item cannot be nil")
	ErrItemNotFound = errors.New("item not found")
)
