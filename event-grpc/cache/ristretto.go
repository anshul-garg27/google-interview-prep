package cache

import (
	"github.com/dgraph-io/ristretto"
	"init.bulbul.tv/bulbul-backend/event-grpc/config"
	"sync"
)

var (
	singletonRistretto *ristretto.Cache
	ristrettoOnce      sync.Once
)

func Ristretto(config config.Config) *ristretto.Cache {
	ristrettoOnce.Do(func() {
		singletonRistretto, _ = ristretto.NewCache(&ristretto.Config{
			NumCounters: 1e7,     // number of keys to track frequency of (10M).
			MaxCost:     1 << 30, // maximum cost of cache (1GB).
			BufferItems: 64,      // number of keys per Get buffer.
		})
	})
	return singletonRistretto
}
