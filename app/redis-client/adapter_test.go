package redisclient_test

import (
	"testing"

	redisclient "github.com/Outtech105k/ShortUrlServer/app/redis-client"
	"github.com/stretchr/testify/assert"
)

func TestAdapter_BasicMethods(t *testing.T) {
	mr, adapter, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Run("Set and Get", func(t *testing.T) {
		key := "test-key"
		val := "test-value"

		err := adapter.Set(key, val)
		assert.NoError(t, err)

		got, err := adapter.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, val, got)
	})

	t.Run("IsExists", func(t *testing.T) {
		key := "exists-key"
		mr.Set(key, "val")

		exists, err := adapter.IsExists(key)
		assert.NoError(t, err)
		assert.True(t, exists)

		exists, err = adapter.IsExists("non-existent")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("IsExists - Error", func(t *testing.T) {
		mr.Close() // 接続を切断してエラーを発生させる
		_, err := adapter.IsExists("key")
		assert.Error(t, err)
	})
}

func TestNewRedisAdapter(t *testing.T) {
	t.Run("Error - Invalid Address", func(t *testing.T) {
		_, err := redisclient.NewRedisAdapter("invalid-addr", "", 0)
		assert.Error(t, err)
	})
}
