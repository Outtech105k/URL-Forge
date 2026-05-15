package redisclient_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetURLRecord(t *testing.T) {
	t.Run("Success without Expiration", func(t *testing.T) {
		mr, adapter, cleanup := setupTestEnvironment(t)
		defer cleanup()

		id := "id11"
		baseURL := "https://example.com/target"

		err := adapter.SetURLRecord(id, baseURL, true, true, nil)
		assert.NoError(t, err)
		assert.Equal(t, baseURL, mr.HGet(id, "base_url"))

		isSandCushion, _ := strconv.ParseBool(mr.HGet(id, "cushion"))
		assert.Equal(t, true, isSandCushion)

		isPublicCtrl, _ := strconv.ParseBool(mr.HGet(id, "public_ctrl"))
		assert.Equal(t, true, isPublicCtrl)
	})

	t.Run("Success with Expiration", func(t *testing.T) {
		mr, adapter, cleanup := setupTestEnvironment(t)
		defer cleanup()

		id := "id12"
		baseURL := "https://example.com/target2"
		expire := 10 * time.Minute

		err := adapter.SetURLRecord(id, baseURL, false, false, &expire)
		assert.NoError(t, err)
		assert.Equal(t, baseURL, mr.HGet(id, "base_url"))
		assert.True(t, mr.TTL(id) > 0)

		isPublicCtrl, _ := strconv.ParseBool(mr.HGet(id, "public_ctrl"))
		assert.Equal(t, false, isPublicCtrl)
	})

	t.Run("Error - HMSet Failure", func(t *testing.T) {
		mr, adapter, _ := setupTestEnvironment(t)
		// cleanup を呼ばずに Close してエラーを誘発
		adapter.Close()
		mr.Close()

		err := adapter.SetURLRecord("id", "url", false, true, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "setRecord")
	})

	t.Run("Error - Expire Failure", func(t *testing.T) {
		mr, adapter, cleanup := setupTestEnvironment(t)
		defer cleanup()

		expire := 10 * time.Minute
		// HMSet は成功させ、その後に Close して Expire を失敗させるのは難しいが
		// インターフェース化していない現状では HMSet 自体がエラーになる
		adapter.Close()
		mr.Close()

		err := adapter.SetURLRecord("id", "url", false, true, &expire)
		assert.Error(t, err)
	})
}
