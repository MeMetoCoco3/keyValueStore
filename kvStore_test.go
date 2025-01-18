package main

import (
	"sync"
	"testing"
)

func TestKVStore(t *testing.T) {
	store1 := NewKVStore[string, int]()

	// Put and Get
	err := store1.Put("key1", 100)
	if err != nil {
		t.Errorf("Put failed: %v", err)
	}

	val, err := store1.Get("key1")
	if err != nil || val != 100 {
		t.Errorf("Get failed: expected 100, got %v, err: %v", val, err)
	}

	// Put Duplicate key
	err = store1.Put("key1", 200)
	if err == nil {
		t.Errorf("Put should have failed for duplicate key")
	}

	// Get Nonexistent Key
	_, err = store1.Get("key2")
	if err == nil {
		t.Errorf("Get should have failed for non-existent key")
	}

	// Update
	err = store1.Update("key1", 300)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}
	val, _ = store1.Get("key1")
	if val != 300 {
		t.Errorf("Update failed: expected 300, got %v", val)
	}

	// Update Nonexistent key
	err = store1.Update("key1", 300)
	if err != nil {
		t.Errorf("We updated non existent key: %v", err)
	}

	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	// Delete
	deletedVal, err := store1.Delete("key1")
	if err != nil || deletedVal != 300 {
		t.Errorf("Delete failed: expected 300, got %v, err: %v", deletedVal, err)
	}

	// Delete Nonexistent key
	_, err = store1.Delete("key1")
	if err == nil {
		t.Errorf("Delete should have failed for non-existent key")
	}

	// Has
	store1.Put("key1", 500)
	if ok := store1.Has("key1"); !ok {
		t.Errorf("Store has 'key1'")
	}

	if ok := store1.Has("We do not have this key"); ok {
		t.Errorf("We do not have this key")
	}

	// Clear
	store1.Clear()
	_, err = store1.Get("key2")
	if err == nil {
		t.Errorf("Clear failed: store1 should be empty")
	}

	// Other Datatypes
	store2 := NewKVStore[uint8, float32]()
	store2.Put(2, 3.14)
	val2, err := store2.Get(2)
	if err != nil {
		t.Errorf("Get [uint8, float32] failed, got %v: %v", val2, err)
	}

	err = store2.Update(2, 6.28)
	if err != nil {
		t.Errorf("Update [uint8, float32] failed: expected 6,28, got %v", val)

	}

	_, err = store2.Delete(2)
	if err != nil {
		t.Errorf("Delete [uint8, float32] failed: %v", val)
	}
}

func TestKVStoreConcurrency(t *testing.T) {
	store1 := NewKVStore[int, int]()
	var wg sync.WaitGroup
	n := 100

	// Concurrent writes
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = store1.Put(i, i*10)
		}(i)
	}
	wg.Wait()

	// Concurrent reads
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			val, err := store1.Get(i)
			if err != nil || val != i*10 {
				t.Errorf("Concurrent Get failed for key %d: got %d, err: %v", i, val, err)
			}
		}(i)
	}
	wg.Wait()
}
