package watchers

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/lordvidex/gostream/internal/entity"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// WatcherRegistrar stores streams for clients and updates them as needed.
type WatcherRegistrar struct {
	mu *sync.RWMutex
	// watchers stores map of watched entities (w) to map of clients (c)
	// client map (c) contains key identifier for client and their corresponding update channel (u)
	watchers map[gostreamv1.Entity]map[*entity.Watcher]struct{}
	count    atomic.Int64
}

// NewWatcherRegistrar ...
func NewWatcherRegistrar() *WatcherRegistrar {
	return &WatcherRegistrar{
		watchers: make(map[gostreamv1.Entity]map[*entity.Watcher]struct{}),
		mu:       new(sync.RWMutex),
	}
}

// RegisterWatcher stores the update channel for a client
func (c *WatcherRegistrar) RegisterWatcher(w *entity.Watcher) error {
	if w == nil {
		return errors.New("nil watcher")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ent := range w.Entities() {
		if _, ok := c.watchers[ent]; !ok {
			c.watchers[ent] = make(map[*entity.Watcher]struct{})
		}
		c.watchers[ent][w] = struct{}{}
	}

	c.count.Add(1)
	// TODO: send snapshot first

	return nil
}

// DeleteWatcher ...
func (c *WatcherRegistrar) DeleteWatcher(w *entity.Watcher) error {
	if w == nil {
		return errors.New("nil watcher")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	w.Close() // safe to close more than once

	found := false
	for _, ent := range w.Entities() {
		if !found {
			if _, ok := c.watchers[ent]; ok {
				found = true
				c.count.Add(-1)
			}
		}
		delete(c.watchers[ent], w)
	}

	return nil
}

// Count is the number of currently listening watchers
func (c *WatcherRegistrar) Count() int64 {
	return c.count.Load()
}

// PublishToClients propagates updates to registered watchers
func (c *WatcherRegistrar) PublishToClients(ctx context.Context, data *gostreamv1.WatchResponse) error {

	var channels []gostreamv1.Entity

	switch data.GetKind() {
	case gostreamv1.EventKind_EVENT_KIND_UPDATE:
		update := data.GetUpdate()
		if update == nil {
			return errors.New("kind is update, but update is nil")
		}
		entity := update.GetEntity()
		channels = []gostreamv1.Entity{entity, gostreamv1.Entity_ENTITY_UNSPECIFIED}
	case gostreamv1.EventKind_EVENT_KIND_DELETE:
		del := data.GetDelete()
		if del == nil {
			return errors.New("kind is delete, but delete is nil")
		}
		entity := del.GetEntity()
		channels = []gostreamv1.Entity{entity, gostreamv1.Entity_ENTITY_UNSPECIFIED}
	}

	if len(channels) == 0 {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, updateChannel := range channels {
		for watcher := range c.watchers[updateChannel] {
			watcher.Send(data)
		}
	}

	return nil
}

// Close removes clients from watchers and closes their update channels
func (c *WatcherRegistrar) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, clients := range c.watchers {
		for client := range clients {
			client.Close()
		}
	}

	clear(c.watchers)
	fmt.Println("client watcher closed")
	return nil
}
