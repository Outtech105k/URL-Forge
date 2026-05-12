package routes_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Outtech105k/ShortUrlServer/app/routes"
	"github.com/Outtech105k/ShortUrlServer/app/testutils"
	"github.com/Outtech105k/ShortUrlServer/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// templates ディレクトリにアクセスできるようにディレクトリを移動
	oldWd, _ := os.Getwd()
	if err := os.Chdir(".."); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer os.Chdir(oldWd)

	mockRedis := new(testutils.MockRedisClient)
	appCtx := &utils.AppContext{
		Config: utils.Config{
			ServerEndpoint: "https://srv.test",
		},
		Redis: mockRedis,
	}

	router := routes.SetupRouter(appCtx)

	t.Run("GET / should render index.html", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "<title>") // index.html の中身があるか
	})

	t.Run("POST /set should reach SetUrlHandler", func(t *testing.T) {
		// ハンドラーまで到達することを確認（バリデーションエラーで400が返ればルートは正しい）
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/set", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GET /:shortUrl should reach GetUrlHandler", func(t *testing.T) {
		mockRedis.On("GetBaseUrl", "test").Return("", nil).Once()
		mockRedis.On("GetIsNeedCusionPage", "test").Return(false, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		// 302 リダイレクトまたは 404/500 ならハンドラーに到達している
		assert.NotEqual(t, http.StatusNotFound, w.Code)
		mockRedis.AssertExpectations(t)
	})

	t.Run("CORS configuration", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/set", nil)
		req.Header.Set("Origin", "http://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	})
}
