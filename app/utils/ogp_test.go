package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetchOGPInfo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `
<html>
<head>
    <meta property="og:title" content="Test Title" />
    <meta property="og:description" content="Test Description" />
    <meta property="og:image" content="https://example.com/image.png" />
</head>
<body>Hello</body>
</html>`)
		}))
		defer ts.Close()

		info, err := FetchOGPInfo(ts.URL, 5*time.Second)
		assert.NoError(t, err)
		assert.Equal(t, "Test Title", info.Title)
		assert.Equal(t, "Test Description", info.Description)
		assert.Equal(t, "https://example.com/image.png", info.Image)
		assert.NotEmpty(t, info.Domain)
	})

	t.Run("FetchError", func(t *testing.T) {
		// 存在しないURL
		info, err := FetchOGPInfo("http://localhost:1", 5*time.Second)
		assert.Error(t, err)
		assert.NotNil(t, info)
		assert.Equal(t, "http://localhost:1", info.URL)
		assert.Equal(t, "localhost:1", info.Domain)
	})

	t.Run("NoOGPTags", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `<html><body>No OGP here</body></html>`)
		}))
		defer ts.Close()

		info, err := FetchOGPInfo(ts.URL, 5*time.Second)
		assert.NoError(t, err)
		assert.Equal(t, "", info.Title)
		assert.Equal(t, "", info.Description)
		assert.Equal(t, "", info.Image)
	})
}
