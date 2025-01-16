package main

import (
	"fmt"
	"sync"
)

// Convention -er means interface.
// Comparable, any...
type Storer[K comparable, V any] interface {
	Put(K, V) error
	Get(K) (V, error)
	Update(K, V) error
	Delete(K) (V, error)
}

type KVStore[K comparable, V any] struct {
	// Allows to bloc
	mu   sync.RWMutex
	data map[K]V
}

func NewKVStore[K comparable, V any]() *KVStore[K, V] {
	return &KVStore[K, V]{
		mu:   sync.RWMutex{},
		data: make(map[K]V),
	}
}

func (s *KVStore[K, V]) Put(key K, val V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if v, ok := s.data[key]; ok {
		return fmt.Errorf("Found Val: %v on Key: %v", v, key)
	}
	s.data[key] = val
	return nil
}

func (s *KVStore[K, V]) Get(key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if v, ok := s.data[key]; !ok {
		// Cuidao, null value
		var v V
		return v, fmt.Errorf("Value not found on key %v", key)
	} else {
		return v, nil
	}
}

func (s *KVStore[K, V]) Update(key K, val V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = val
	return nil
}

// Delete and clear cool
func (s *KVStore[K, V]) Delete(key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, err := s.Get(key)
	if err != nil {
		var v V
		return v, fmt.Errorf("Value not found on key %v", key)
	}

	delete(s.data, key)
	return v, nil
}

func (s *KVStore[K, V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.data)
}
