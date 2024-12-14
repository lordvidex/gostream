// Package entity holds structs shared between packages
package entity

import (
	"slices"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// WatcherOpts ...
type WatcherOpts func(*Watcher)

// WithWatcherIdentifier ...
func WithWatcherIdentifier(identifier string) WatcherOpts {
	return func(w *Watcher) {
		w.identifier = identifier
	}
}

// Watcher ...
type Watcher struct {
	identifier string
	created    time.Time
	lastSent   time.Time
	updates    chan<- *gostreamv1.WatchResponse
	entities   []gostreamv1.Entity
	closed     atomic.Bool
}

// NewWatcher ...
func NewWatcher(updateChan chan<- *gostreamv1.WatchResponse, kind []gostreamv1.Entity, opts ...WatcherOpts) *Watcher {
	w := Watcher{
		identifier: uuid.New().String(),
		created:    time.Now(),
		updates:    updateChan,
		entities:   collapse(kind),
	}

	for _, opt := range opts {
		opt(&w)
	}
	return &w
}

// Close should be called to clean up client resources.
// Close is safe to call more than once.
func (w *Watcher) Close() {
	if w.closed.Load() {
		return
	}
	w.closed.Store(true)
	close(w.updates)
}

// IsClosed ...
func (w *Watcher) IsClosed() bool {
	return w.closed.Load()
}

// Send ...
func (w *Watcher) Send(update *gostreamv1.WatchResponse) {
	if w.closed.Load() {
		return
	}
	w.updates <- update
}

// Entities ...
func (w *Watcher) Entities() []gostreamv1.Entity {
	return w.entities
}

func collapse(entities []gostreamv1.Entity) []gostreamv1.Entity {
	entities = slices.Compact(entities)
	if slices.Contains(entities, gostreamv1.Entity_ENTITY_UNSPECIFIED) {
		return []gostreamv1.Entity{gostreamv1.Entity_ENTITY_UNSPECIFIED}
	}

	return entities
}
