package inmemory

import (
	"iter"
	"sync"
)

// Array is a Container, that's safe for concurrent access and has
// 1. Fast lookups O(1)
// 2. Fast inserts O(1)
// 3. Slow delete O(n)
type Array[K comparable, V Value[K]] struct {
	pos  map[K]int
	data []V
	mu   sync.RWMutex
}

func NewArray[K comparable, V Value[K]]() *Array[K, V] {
	return &Array[K, V]{
		data: make([]V, 0),
		pos:  make(map[K]int),
	}
}

func (a *Array[K, V]) Add(v V) {
	a.mu.Lock()
	a.add(v)
	a.mu.Unlock()
}

func (a *Array[K, V]) add(v V) {
	k := v.Key()
	if idx, ok := a.pos[k]; ok {
		a.data[idx] = v
	} else {
		a.pos[k] = len(a.data)
		a.data = append(a.data, v)
	}
}

func (a *Array[K, V]) AddAll(vs []V) {
	a.mu.Lock()
	for _, v := range vs {
		a.add(v)
	}
	a.mu.Unlock()
}

// Remove has a time complexity of O(n)
func (a *Array[K, V]) Remove(k K) {
	a.mu.Lock()
	if idx, ok := a.pos[k]; ok {
		delete(a.pos, k)
		a.data = append(a.data[:idx], a.data[idx+1:]...)
		for i := idx; i < len(a.data); i++ {
			a.pos[a.data[i].Key()] = i
		}
	}
	a.mu.Unlock()
}

// Snapshot simply returns data array
func (a *Array[K, V]) Snapshot() []V {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.data
}

func (a *Array[K, V]) Get(k K) (V, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if idx, ok := a.pos[k]; ok {
		return a.data[idx], true
	}

	var v V // alternative to *new(V) but without alloc
	return v, false
}

func (a *Array[K, V]) Clear() {
	a.mu.Lock()
	a.pos = make(map[K]int)
	a.data = make([]V, 0)
	a.mu.Unlock()
}

func (a *Array[K, V]) Iter() iter.Seq[V] {
	return func(yield func(V) bool) {
		a.mu.RLock()
		defer a.mu.RUnlock()

		for _, v := range a.data {
			if !yield(v) {
				return
			}
		}
	}
}
