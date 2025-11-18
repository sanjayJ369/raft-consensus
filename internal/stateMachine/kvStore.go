package statemachine

import (
	"fmt"

	"github.com/sanjayJ369/raft-consensus/internal/types"
)

// KVStore is a simple in memory
// key value store with simple interface
// supports: Get, Set, Del
type KVStore struct {
	store map[string]any
	lgr   types.Logger
}

func NewKVStore(lgr types.Logger) types.DB {
	return &KVStore{
		store: make(map[string]any),
		lgr:   lgr,
	}
}

// Get returns value of it exists
// or returns nil if it does not exist
func (k *KVStore) Get(key string) any {
	k.lgr.Logf("get key: %s", key)
	if val, ok := k.store[key]; ok {
		return val
	}
	return nil
}

// Put sets key => val
// returns if the operation was successful
func (k *KVStore) Put(key string, val any) error {
	k.lgr.Logf("put key %s: val %v", key, val)
	k.store[key] = val
	return nil
}

// Del delets the key if it exists
// returns true on scuessful deletion
// if the element does not exists returns true by default
func (k *KVStore) Del(key string) error {
	k.lgr.Logf("del key: %s", key)
	if _, ok := k.store[key]; !ok {
		return fmt.Errorf("element not found: %s", key)
	}
	delete(k.store, key)
	return nil
}
