package controllers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControlUrlHandler(t *testing.T) {
	_, mr, router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// OGP情報を提供するモックサーバー
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `<html><head><title>Mock Title</title></head></html>`)
	}))
	defer ts.Close()

	t.Run("Success", func(t *testing.T) {
		shortUrl := "testshort"
		mr.HSet(shortUrl, "base_url", ts.URL)
		mr.HSet(shortUrl, "public_ctrl", "true")

		w := performRequest(router, "GET", "/"+shortUrl+"/control", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Mock Title")
		assert.Contains(t, w.Body.String(), ts.URL)
	})

	t.Run("Forbidden (Private Control)", func(t *testing.T) {
		shortUrl := "privateshort"
		mr.HSet(shortUrl, "base_url", "https://example.com")
		mr.HSet(shortUrl, "public_ctrl", "false")

		w := performRequest(router, "GET", "/"+shortUrl+"/control", nil)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "このURLの管理画面は非公開です。")
	})

	t.Run("NotFound", func(t *testing.T) {
		w := performRequest(router, "GET", "/nonexistent/control", nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
