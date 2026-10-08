# go-cache

In-memory cache abstractions for Go: a small `Cache` interface and a thread-safe cache whose items expire at a given time.

## Installation

```bash
go get github.com/ralvarezdev/go-cache
```

## Usage

```go
import (
	"fmt"
	"time"

	gocachetimed "github.com/ralvarezdev/go-cache/timed"
)

cache := gocachetimed.NewDefaultTimedCache()

// Values must be wrapped in a *TimedItem with an expiration time
item := gocachetimed.NewTimedItem("hello", time.Now().Add(5*time.Minute))
if err := cache.Set("greeting", item); err != nil {
	panic(err)
}

if value, ok := cache.Get("greeting"); ok {
	fmt.Println(value) // "hello"
}

cache.Delete("greeting")
```

## API

**Root package `gocache`**

- **`Cache`** — interface: `Set(key, value) error`, `UpdateValue(key, value) error`, `Has(key) bool`, `Get(key) (interface{}, bool)`, `Delete(key)`.
- **Errors** — `ErrNilCache`, `ErrNilItem`, `ErrItemNotFound`.

**Package `timed`**

- **`TimedCache`** — embeds `gocache.Cache`, adds `GetExpirationTime(key)` and `UpdateExpirationTime(key, time.Time)`.
- **`TimedItem`** — `NewTimedItem(value, expiresAt)`, `GetValue`, `SetValue`, `GetExpiresAt`, `SetExpiresAt`, `HasExpired`.
- **`DefaultTimedCache`** — mutex-protected map from `NewDefaultTimedCache()`. `Set` requires a `*TimedItem` that has not expired; `Get` returns `(nil, false)` for missing or expired items and removes expired ones. Also has `UpdateExpiresAt`.
- **Errors** — `ErrItemHasExpired`, `ErrValueMustBeATimedItem`.

There are no tests.

## License

GNU General Public License v3.0. See [LICENSE](LICENSE).
