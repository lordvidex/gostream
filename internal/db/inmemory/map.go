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
	//TODO implement me
	panic("implement me")
}

func (m *Map[K, V]) AddAll(vs []V) {
	//TODO implement me
	panic("implement me")
}

func (m *Map[K, V]) Remove(k K) {
	//TODO implement me
	panic("implement me")
}

func (m *Map[K, V]) Get(k K) (V, bool) {
	//TODO implement me
	panic("implement me")
}

func (m *Map[K, V]) Clear() {
	//TODO implement me
	panic("implement me")
}

func (m *Map[K, V]) Snapshot() []V {
	//TODO implement me
	panic("implement me")
}

func (m *Map[K, V]) Iter() iter.Seq[V] {
	//TODO implement me
	panic("implement me")
}
