package inmemory

import (
	"container/list"
	"iter"
	"sync"
)

// Map is a Container that's safe for concurrent use and has:
// 1. Fast lookups
// 2. Fast inserts
// 3. Fast delete
// 4. Snapshot is O(n)
type Map[K comparable, V Value[K]] struct {
	data map[K]*list.Element
	l    list.List
	mu   sync.RWMutex
}

func (m *Map[K, V]) lazyInit() {
	if m.data == nil {
		m.data = make(map[K]*list.Element)
	}
}

func NewMap[K comparable, V Value[K]]() *Map[K, V] {
	return &Map[K, V]{
		data: make(map[K]*list.Element),
		l:    list.List{},
	}
}

func (m *Map[K, V]) Add(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lazyInit()
	el := m.l.PushBack(v)
	m.data[k] = el

}

func (m *Map[K, V]) AddAll(vs []V) {
	m.mu.Lock()
	m.lazyInit()
	for _, v := range vs {
		el := m.l.PushBack(v)
		m.data[v.Key()] = el
	}
	m.mu.Unlock()
}

func (m *Map[K, V]) Remove(k K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	el := m.data[k]
	if el != nil {
		m.l.Remove(el)
	}
	delete(m.data, k)
}

func (m *Map[K, V]) Get(k K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.data[k]
	if !ok {
		var v V
		return v, false
	}
	return value.Value.(V), true
}

func (m *Map[K, V]) Clear() {
	m.mu.Lock()
	m.data = make(map[K]*list.Element)
	m.l = list.List{}
	m.mu.Unlock()
}

func (m *Map[K, V]) Snapshot() []V {
	m.mu.RLock()
	values := make([]V, 0, len(m.data))
	m.mu.RUnlock()
	for value := range m.Iter() {
		values = append(values, value)
	}
	return values
}

func (m *Map[K, V]) Iter() iter.Seq[V] {
	return func(yield func(V) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		for el := m.l.Front(); el != nil; el = el.Next() {
			if !yield(el.Value.(V)) {
				return
			}
		}

	}
}
