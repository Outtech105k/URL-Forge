package models_test

import (
	"encoding/json"
	"testing"

	"github.com/Outtech105k/ShortUrlServer/app/models"
	"github.com/stretchr/testify/assert"
)

func TestJSON_Serialization(t *testing.T) {
	t.Run("APIError JSON tags", func(t *testing.T) {
		err := models.APIError{
			Type:    "error_type",
			Message: "error message",
		}
		data, _ := json.Marshal(err)

		var m map[string]interface{}
		json.Unmarshal(data, &m)

		assert.Equal(t, "error_type", m["type"])
		assert.Equal(t, "error message", m["message"])
		assert.Nil(t, m["details"]) // omitempty の確認
	})

	t.Run("APIResponse JSON tags", func(t *testing.T) {
		resp := models.APIResponse{
			BaseURL:  "https://example.com",
			ShortURL: "https://srv.test/id",
		}
		data, _ := json.Marshal(resp)

		var m map[string]interface{}
		json.Unmarshal(data, &m)

		assert.Equal(t, "https://example.com", m["base_url"])
		assert.Equal(t, "https://srv.test/id", m["short_url"])
	})

	t.Run("SetUrlRequest JSON tags", func(t *testing.T) {
		// すべてのフィールドが指定された場合
		customID := "custom"
		useTrue := true
		var idLen uint32 = 8

		req := models.SetUrlRequest{
			BaseURL:      "https://example.com",
			CustomID:     &customID,
			UseUppercase: &useTrue,
			IDLength:     &idLen,
		}
		data, _ := json.Marshal(req)

		var m map[string]interface{}
		json.Unmarshal(data, &m)

		assert.Equal(t, "https://example.com", m["base_url"])
		assert.Equal(t, "custom", m["custom_id"])
		assert.Equal(t, true, m["use_uppercase"])
		assert.Equal(t, float64(8), m["id_length"])
		assert.Nil(t, m["use_lowercase"]) // 指定しなかったフィールドが nil (JSON では omit される)
	})
}
