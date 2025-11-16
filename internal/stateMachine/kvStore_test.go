package statemachine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKVStore(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		store := KVStore{
			store: make(map[string]any),
		}

		store.Put("hi", "val")
		want := "val"
		got := store.Get("hi")

		assert.Equal(t, want, got)

		val := store.Get("no key")
		assert.Nil(t, val)
	})

	t.Run("del", func(t *testing.T) {
		store := KVStore{
			store: make(map[string]any),
		}

		key := "hi"
		val := "hello"

		store.Put(key, val)
		got := store.Get(key)
		assert.Equal(t, val, got)

		err := store.Del("hi")
		require.NoError(t, err)

		err = store.Del("hi")
		require.Error(t, err)
	})

	t.Run("put", func(t *testing.T) {
		k1 := "hello"
		v1 := "world"

		store := KVStore{
			store: make(map[string]any),
		}

		err := store.Put(k1, v1)
		require.NoError(t, err)
		got := store.Get(k1)
		assert.Equal(t, v1, got)

		v2 := "universe"
		err = store.Put(k1, v2)
		require.NoError(t, err)
		got = store.Get(k1)
		assert.Equal(t, v2, got)
	})
}
