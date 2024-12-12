package inmemory

import (
	"iter"
	"sync"
)

// Map is a Container that's safe for concurrent use and has:
// 1. Fast lookups
// 2. Fast inserts
// 3. Fast delete
// 4. Snapshot requires sorting
type Map[K comparable, V Value[K]] struct {
	data map[K]V
	mu   sync.RWMutex
}

func (m *Map[K, V]) Add(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[k] = v
}

func (m *Map[K, V]) AddAll(vs []V) {
	m.mu.Lock()
	for _, v := range vs {
		m.data[v.Key()] = v
	}
	m.mu.Unlock()
}

func (m *Map[K, V]) Remove(k K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data != nil {
		delete(m.data, k)
	}
}

func (m *Map[K, V]) Get(k K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for keys := range m.data {
		if keys == k {
			return m.data[keys], true
		}
	}
	var v V // alternative to *new(V) but without alloc
	return v, false
}

func (m *Map[K, V]) Clear() {
	m.mu.Lock()
	m.data = make(map[K]V)
	m.mu.Unlock()
}

func (m *Map[K, V]) Snapshot() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	values := make([]V, 0, len(m.data))
	for _, value := range m.data {
		values = append(values, value)
	}
	return values
}

func (m *Map[K, V]) Iter() iter.Seq[V] {
	return func(yield func(V) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		for _, v := range m.data {
			if !yield(v) {
				return
			}
		}
	}
}
