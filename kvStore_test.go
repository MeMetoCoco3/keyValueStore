package main

import (
	"sync"
	"testing"
)

func TestKVStore(t *testing.T) {
	store := NewKVStore[string, int]()

	err := store.Put("key1", 100)
	if err != nil {
		t.Errorf("Put failed: %v", err)
	}

	// Duplicate key
	err = store.Put("key1", 200)
	if err == nil {
		t.Errorf("Put should have failed for duplicate key")
	}

	val, err := store.Get("key1")
	if err != nil || val != 100 {
		t.Errorf("Get failed: expected 100, got %v, err: %v", val, err)
	}

	// Nonexistent Key
	_, err = store.Get("key2")
	if err == nil {
		t.Errorf("Get should have failed for non-existent key")
	}

	err = store.Update("key1", 300)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}
	val, _ = store.Get("key1")
	if val != 300 {
		t.Errorf("Update failed: expected 300, got %v", val)
	}

	// Test Delete
	deletedVal, err := store.Delete("key1")
	if err != nil || deletedVal != 300 {
		t.Errorf("Delete failed: expected 300, got %v, err: %v", deletedVal, err)
	}

	// Nonexistent key
	_, err = store.Delete("key1")
	if err == nil {
		t.Errorf("Delete should have failed for non-existent key")
	}

	store.Put("key2", 500)
	store.Clear()
	_, err = store.Get("key2")
	if err == nil {
		t.Errorf("Clear failed: store should be empty")
	}
}

func TestKVStoreConcurrency(t *testing.T) {
	store := NewKVStore[int, int]()
	var wg sync.WaitGroup
	n := 100

	// Concurrent writes
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = store.Put(i, i*10)
		}(i)
	}
	wg.Wait()

	// Concurrent reads
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			val, err := store.Get(i)
			if err != nil || val != i*10 {
				t.Errorf("Concurrent Get failed for key %d: got %d, err: %v", i, val, err)
			}
		}(i)
	}
	wg.Wait()
}
