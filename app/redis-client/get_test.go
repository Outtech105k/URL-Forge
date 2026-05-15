package redisclient_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBaseUrl(t *testing.T) {
	mr, adapter, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Run("Success", func(t *testing.T) {
		key := "id1"
		baseUrl := "https://example.com/target"
		mr.HSet(key, "base_url", baseUrl)

		result, err := adapter.GetBaseUrl(key)
		assert.NoError(t, err)
		assert.Equal(t, baseUrl, result)
	})

	t.Run("Not Found", func(t *testing.T) {
		_, err := adapter.GetBaseUrl("non-existent")
		assert.Error(t, err)
	})
}

func TestGetIsNeedCusionPage(t *testing.T) {
	mr, adapter, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Run("Success - True", func(t *testing.T) {
		key := "id2"
		mr.HSet(key, "cushion", "true")

		result, err := adapter.GetIsNeedCusionPage(key)
		assert.NoError(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("Success - False", func(t *testing.T) {
		key := "id3"
		mr.HSet(key, "cushion", "false")

		result, err := adapter.GetIsNeedCusionPage(key)
		assert.NoError(t, err)
		assert.Equal(t, false, result)
	})

	t.Run("Error - Parse Failure", func(t *testing.T) {
		key := "id4"
		mr.HSet(key, "cushion", "invalid-bool")

		_, err := adapter.GetIsNeedCusionPage(key)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parse got val")
	})

	t.Run("Error - Redis Error (Key Missing)", func(t *testing.T) {
		_, err := adapter.GetIsNeedCusionPage("non-existent")
		assert.Error(t, err)
	})
}

func TestGetIsPublicCtrl(t *testing.T) {
	mr, adapter, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Run("Success - True", func(t *testing.T) {
		key := "id_p1"
		mr.HSet(key, "public_ctrl", "true")

		result, err := adapter.GetIsPublicCtrl(key)
		assert.NoError(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("Success - False", func(t *testing.T) {
		key := "id_p2"
		mr.HSet(key, "public_ctrl", "false")

		result, err := adapter.GetIsPublicCtrl(key)
		assert.NoError(t, err)
		assert.Equal(t, false, result)
	})

	t.Run("Error - Parse Failure", func(t *testing.T) {
		key := "id_p3"
		mr.HSet(key, "public_ctrl", "invalid-bool")

		_, err := adapter.GetIsPublicCtrl(key)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parse got val")
	})

	t.Run("Error - Redis Error (Key Missing)", func(t *testing.T) {
		_, err := adapter.GetIsPublicCtrl("non-existent")
		assert.Error(t, err)
	})
}
