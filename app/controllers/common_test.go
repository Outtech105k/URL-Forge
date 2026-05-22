package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	redisclient "github.com/Outtech105k/ShortUrlServer/app/redis-client"
	"github.com/Outtech105k/ShortUrlServer/app/routes"
	"github.com/Outtech105k/ShortUrlServer/app/utils"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
)

// テスト用の共通環境を構築
// 戻り値: コンテキスト, miniredisインスタンス, ルーター, クリーンアップ関数
func setupTestEnvironment(t *testing.T) (utils.AppContext, *miniredis.Miniredis, *gin.Engine, func()) {
	t.Helper()

	// templatesディレクトリへのパス調整 (app/controllers から app/ へ移動)
	oldWd, _ := os.Getwd()
	if err := os.Chdir(".."); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to run miniredis: %v", err)
	}

	adapter, err := redisclient.NewRedisAdapter(mr.Addr(), "", 0)
	if err != nil {
		t.Fatalf("failed to create redis adapter: %v", err)
	}

	appCtx := &utils.AppContext{
		Config: utils.Config{
			ServerEndpoint:  "https://srv.test",
			AppName:         "URL Forge",
			MaxIDLength:     100,
			DefaultIDLength: 6,
			OGPFetchTimeout: 5 * time.Second,
			AllowOrigins:    "*",
			BotUserAgents:   []string{"bot", "crawler", "spider", "facebookexternalhit", "twitterbot", "slackbot", "discordbot", "whatsapp", "line-poker"},
			MaxRetryCount:   10,
		},
		Redis: adapter,
	}

	router := routes.SetupRouter(appCtx)

	cleanup := func() {
		adapter.Close()
		mr.Close()
		if err := os.Chdir(oldWd); err != nil {
			t.Logf("Warning: failed to restore directory: %v", err)
		}
	}

	return *appCtx, mr, router, cleanup
}

// HTTPリクエストを実行してレスポンスを返す
func performRequest(router *gin.Engine, method, path string, headers map[string]string, body interface{}) *httptest.ResponseRecorder {
	var buf *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewBuffer(b)
	} else {
		buf = bytes.NewBuffer(nil)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, buf)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	router.ServeHTTP(w, req)
	return w
}
