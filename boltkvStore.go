package main

import (
	"encoding/json"
	"fmt"
	"sync"

	bolt "go.etcd.io/bbolt"
)

type BoltStorer[K comparable, V any] interface {
	Storer[K, V]
	PutB(K, V) error
	GetB(K) (V, error)
	GetAll() (map[K]V, error)
}

type BoltStore[K comparable, V any] struct {
	*KVStore[K, V] // Struct embedding, inherits methods
	mu             sync.RWMutex
	data           map[K]V
	db             *bolt.DB
	bucket         *bolt.Bucket //Bucket name
}

/*
	type Storer[K comparable, V any] interface {
		Put(K, V) error
		Get(K) (V, error)
		Update(K, V) error
		Delete(K) (V, error)
		Iter(...K) ([]V, error)
	}

	func (s *BoltStore[K, V]) Put(key K, val V) error {
		fmt.Println(s)
		s.mu.Lock()
		fmt.Println("PresegmentationFault")

		defer s.mu.Unlock()

		if v, ok := s.data[key]; ok {
			return fmt.Errorf("Found Val: %v on Key: %v", v, key)
		}
		s.data[key] = val
		return nil
	}
*/
func NewBoltStore[K comparable, V any](path string) (*BoltStore[K, V], error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("Got error opening DB: %s", err)
	}

	var bucket *bolt.Bucket
	if err = db.Update(func(tx *bolt.Tx) error {
		bucket, err = tx.CreateBucketIfNotExists([]byte("Bunny"))
		if err != nil {
			return fmt.Errorf("Error creating bucket: %s \n", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &BoltStore[K, V]{
		mu:     sync.RWMutex{},
		data:   make(map[K]V),
		db:     db,
		bucket: bucket,
	}, nil
}

func (s *BoltStore[K, V]) Put(key K, val V) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	if v, ok := s.data[key]; ok {
		return fmt.Errorf("Found Val: %v on Key: %v", v, key)
	}
	s.data[key] = val
	return nil
}

func (s *BoltStore[K, V]) Get(key K) (V, error) {
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

func (b *BoltStore[K, V]) PutB(key K, val V) error {
	err := b.Put(key, val)
	if err != nil {
		return fmt.Errorf("Error using Put from Boltstore:, %s", err)
	}

	keyBytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("Error converting key %s into json format:, %s", key, err)
	}

	valueBytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("Error converting value %s into json format:, %s", val, err)
	}

	if err = b.db.Update(func(tx *bolt.Tx) error {
		err = b.bucket.Put(keyBytes, valueBytes)
		if err != nil {
			return fmt.Errorf("Error putting Key:%s Val: %s in bucket. %s", key, val, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("Error putting Key:%s Val: %s in bucket.", key, val)
	}
	return nil
}

func (b *BoltStore[K, V]) GetB(key K) (V, error) {
	var v V
	var val []byte

	keyBytes, err := json.Marshal(key)
	if err != nil {
		return v, fmt.Errorf("Error marshalling key '%s' on GetB. %s", key, err)
	}

	if err = b.db.View(func(tx *bolt.Tx) error {
		val = b.bucket.Get(keyBytes)
		return nil
	}); err != nil {
		return v, fmt.Errorf("Error on get transaction. %s", err)
	}

	err = json.Unmarshal(val, v)
	if err != nil {
		return v, fmt.Errorf("Err1or on get transaction. %s, %s", err, val)
	}

	return v, nil

}

func (b *BoltStore[K, V]) GetAll() (map[K]V, error) {
	result := make(map[K]V)

	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Bunny"))
		if bucket == nil {
			return fmt.Errorf("Bucket not found")
		}

		return bucket.ForEach(func(k, v []byte) error {
			var value V
			if err := json.Unmarshal(v, &value); err != nil {
				return fmt.Errorf("Error unmarshaling value: %v", err)
			}

			var key K
			if err := json.Unmarshal(k, &key); err != nil {
				return fmt.Errorf("Error unmarshaling key: %v", err)
			}

			result[key] = value
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
