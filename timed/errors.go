package timed

import (
	"errors"
)

var (
	ErrItemHasExpired        = errors.New("item has expired")
	ErrValueMustBeATimedItem = errors.New("value must be a TimedItem")
)
