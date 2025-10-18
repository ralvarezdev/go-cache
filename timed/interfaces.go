package timed

import (
	"time"

	gocache "github.com/ralvarezdev/go-cache"
)

type (
	// TimedCache is the interface for timed cache implementations
	TimedCache interface {
		gocache.Cache
		GetExpirationTime(key string) time.Time
		UpdateExpirationTime(key string, expirationTime time.Time) error
	}
)
