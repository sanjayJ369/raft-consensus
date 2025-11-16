package statemachine

// KVStore is a simple in memory
// key value store with simple interface
type KVStore struct {
	store map[string]any
}
