package inmemory

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"iter"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lordvidex/errs/v2"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

const (
	validationFrequency = time.Minute * 5
)

// UpdateMode specifies how data is sent to the provided callback when
// cache has been invalidated and fresh data is returned from the datasource.
type UpdateMode int8

const (
	Snapshot UpdateMode = iota
	Diff
)

// Container stores local data
type Container[K comparable, V any] interface {
	Add(K, V)
	AddAll([]V)
	Remove(K)
	Get(K) (V, bool)
	Clear()
	Snapshot() []V
	Iter() iter.Seq[V]
}

// Value ...
type Value[K comparable] interface {
	// Key represents the identifier for this value.
	Key() K
	// Hash returns a string value that should be always UNIQUE if any field of
	// Value changed.
	Hash() string
	// UniqueString is used for Hash
	UniqueString() string
}

// DataSource represents the cold storage.
type DataSource[K comparable, V Value[K]] interface {
	// Hash returns a unique value for each state of the DataSource,
	// Adding, Deleting or Updating any of the fields in the DataSource must change this Hash value.
	Hash(context.Context) (string, error)
	// FetchAll retrieves all the data from DataSource
	FetchAll(context.Context) ([]V, error)
	// Fetch retrieves item for keys from DataSource
	Fetch(context.Context, ...K) ([]V, error)
}

// Cache ...
type Cache[K comparable, V Value[K]] struct {
	container  Container[K, V]
	dataHash   string
	dirty      atomic.Bool
	dataSource DataSource[K, V]
	update     func(gostreamv1.EventKind, []V)
	mode       UpdateMode
	mu         sync.RWMutex
}

type CacheOpts[K comparable, V Value[K]] func(*Cache[K, V])

// WithDataSource adds data source to cache and enables smart cache features
// such as cache invalidation and reloading.
func WithDataSource[K comparable, V Value[K]](ds DataSource[K, V]) CacheOpts[K, V] {
	return func(c *Cache[K, V]) {
		c.dataSource = ds
	}
}

// WithDataSourceUpdateCallback registers a callback that will be invoked
// with the list of differences between the cache and datasource states.
func WithDataSourceUpdateCallback[K comparable, V Value[K]](mode UpdateMode, cb func(gostreamv1.EventKind, []V)) CacheOpts[K, V] {
	return func(c *Cache[K, V]) {
		c.update = cb
		c.mode = mode
	}
}

// NewCache creates a new cache instance for storing V
func NewCache[K comparable, V Value[K]](
	ctx context.Context,
	container Container[K, V],
	opts ...CacheOpts[K, V],
) (*Cache[K, V], error) {

	ch := Cache[K, V]{
		container: container,
	}

	for _, opt := range opts {
		opt(&ch)
	}

	if err := ch.reset(ctx); err != nil {
		return nil, errs.WrapMsg(err, "failed to initialize cache")
	}

	if ch.dataSource != nil {
		go ch.startDataValidation(ctx)
	}

	return &ch, nil
}

// reset discards all data in cache and reloads them from
// datasource.
func (c *Cache[K, V]) reset(ctx context.Context) error {
	if c.dataSource == nil {
		c.container.Clear()
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	data, err := c.dataSource.FetchAll(ctx)
	if err != nil {
		return errs.WrapMsg(err, "failed to fetch cache data from datasource")
	}
	c.dataHash, err = c.dataSource.Hash(ctx)
	if err != nil {
		return errs.WrapMsg(err, "failed to hash cache data from datasource")
	}

	oldData := c.container.Snapshot()

	if c.update != nil {
		switch c.mode {
		case Snapshot:
			go c.update(gostreamv1.EventKind_EVENT_KIND_SNAPSHOT, data)
		case Diff:
			updates, deletes := computeDifferences(oldData, data)
			go func() {
				c.update(gostreamv1.EventKind_EVENT_KIND_UPDATE, updates)
				c.update(gostreamv1.EventKind_EVENT_KIND_DELETE, deletes)
			}()
		}
	}
	c.container.Clear()
	c.container.AddAll(data)

	return nil
}

func (c *Cache[K, V]) startDataValidation(ctx context.Context) {
	clock := time.NewTicker(validationFrequency)
	defer clock.Stop()

	for {
		select {
		case <-clock.C:
			if c.dirty.Load() {
				c.computeHash()
			}
			if err := c.compareWithDataSource(ctx); err != nil {
				fmt.Println("failed to compare with cache data", err)
				continue
			}
			c.dirty.Store(false)
		case <-ctx.Done():
			return
		}
	}
}

func (c *Cache[K, V]) compareWithDataSource(ctx context.Context) error {
	v, err := c.dataSource.Hash(ctx)
	if err != nil {
		return err
	}

	if v == c.dataHash {
		fmt.Println("cache is up to date")
		return nil
	}

	fmt.Println("cache is dirty, resetting cache")
	return c.reset(ctx)
}

func (c *Cache[K, V]) computeHash() {
	prev := ""
	for v := range c.container.Iter() {
		h := md5.New()
		_, _ = io.WriteString(h, prev)
		_, _ = io.WriteString(h, v.UniqueString())
		prev = fmt.Sprintf("%x", h.Sum(nil))
	}
	c.mu.Lock()
	c.dataHash = prev
	c.mu.Unlock()
}

// Store ...
func (c *Cache[K, V]) Store(key K, value V) {
	c.container.Add(key, value)
	c.dirty.Store(true)
}

// Delete is expensive...
func (c *Cache[K, V]) Delete(key K) {
	c.container.Remove(key)
	c.dirty.Store(true)
}

// Get ...
func (c *Cache[K, V]) Get(key K) (V, bool) {
	return c.container.Get(key)
}

// Snapshot returns an unordered list of items in the cache.
func (c *Cache[K, V]) Snapshot() []V {
	return c.container.Snapshot()
}

func computeDifferences[K comparable, V Value[K]](initial, newData []V) (updates, deletes []V) {
	exist := make(map[K]string)
	for _, v := range initial {
		exist[v.Key()] = v.Hash()
	}
	for _, v := range newData {
		newHash := v.Hash()
		if oldHash, ok := exist[v.Key()]; ok {
			delete(exist, v.Key())
			if newHash == oldHash {
				continue
			}
			updates = append(updates, v)
			continue
		}
		updates = append(updates, v)
	}

	for _, v := range initial {
		if _, ok := exist[v.Key()]; ok {
			deletes = append(deletes, v)
		}
	}
	return updates, deletes
}
