package timed

import (
	"errors"
)

var (
	ErrExpiredItem = errors.New("item has expired")
)
