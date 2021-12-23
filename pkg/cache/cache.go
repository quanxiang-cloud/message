package cache

import (
	"flag"
	"github.com/allegro/bigcache/v3"
	"time"
)

var (
	// number of shards (must be a power of 2)
	shards int
	// time after which entry can be evicted, value is minute
	expired int
	// Interval between removing expired entries, value is second
	clearTime int
	// max entry size in bytes
	maxEntrySize int
	// cache will not allocate more memory than this limit, value in MB
	maxCacheSize int
	//  rps * lifeWindow, used only in initial memory allocation
	maxEntriesInWindow int
)

// PrepareCache PrepareCache
func PrepareCache() {
	flag.IntVar(&shards, "shards", 1024, "number of cache shards")
	flag.IntVar(&expired, "expired", 10, "time after which entry can be evicted, value is minute")
	flag.IntVar(&clearTime, "clear-time", 10, "interval between removing expired entries, value is second")
	flag.IntVar(&maxEntrySize, "max-entry-size", 500, "max entry size in bytes")
	flag.IntVar(&maxCacheSize, "max-cache-size", 500, "cache will not allocate more memory than this limit, value in MB")
	flag.IntVar(&maxEntriesInWindow, "max-entries-in-window", 1000*10*60, "rps * lifeWindow")
}

// NewCache NewCache
func NewCache() (Cache, error) {
	bigCache, err := bigcache.NewBigCache(bigcache.Config{
		Shards:             shards,
		LifeWindow:         time.Duration(expired) * time.Minute,
		CleanWindow:        time.Duration(clearTime) * time.Second,
		MaxEntriesInWindow: maxEntriesInWindow,
		MaxEntrySize:       maxEntrySize,
		StatsEnabled:       false,
		Verbose:            false,
		HardMaxCacheSize:   maxCacheSize,
	})
	if err != nil {
		return nil, err
	}
	return &cache{
		bigCache: bigCache,
	}, nil
}

// Cache Cache
type Cache interface {
	Get(key string) ([]byte, error)
	Push(key string, value []byte) error
}

type cache struct {
	bigCache *bigcache.BigCache
}

func (c *cache) Get(key string) ([]byte, error) {
	return c.bigCache.Get(key)
}

func (c *cache) Push(key string, value []byte) error {
	return c.bigCache.Set(key, value)
}
