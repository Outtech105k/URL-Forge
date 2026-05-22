package controllers_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/Outtech105k/ShortUrlServer/app/controllers"
	"github.com/Outtech105k/ShortUrlServer/app/testutils"
	"github.com/Outtech105k/ShortUrlServer/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestGetUrlHandler_TableDriven(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		shortUrl       string
		userAgent      string
		setupMock      func(m *testutils.MockRedisClient)
		expectedStatus int
		verifyResponse func(t *testing.T, w *http.Response, body string)
	}{
		{
			name:     "Success - Redirect",
			shortUrl: "valid-id",
			setupMock: func(m *testutils.MockRedisClient) {
				m.On("GetBaseUrl", "valid-id").Return("https://example.com", nil).Once()
				m.On("GetIsNeedCusionPage", "valid-id").Return(false, nil).Once()
			},
			expectedStatus: http.StatusFound,
			verifyResponse: func(t *testing.T, w *http.Response, body string) {
				assert.Equal(t, "https://example.com", w.Header.Get("Location"))
			},
		},
		{
			name:     "Success - Show Cushion",
			shortUrl: "cushion-id",
			setupMock: func(m *testutils.MockRedisClient) {
				m.On("GetBaseUrl", "cushion-id").Return("https://example.com", nil).Once()
				m.On("GetIsNeedCusionPage", "cushion-id").Return(true, nil).Once()
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, w *http.Response, body string) {
				assert.Contains(t, body, "https://example.com")
				// OGPタグの確認（クッションページ用）
				assert.Contains(t, body, "<meta property=\"og:title\" content=\"リンクの確認 - URL Forge\" />")
				assert.Contains(t, body, "<meta property=\"og:description\" content=\"この先は外部サイトです。リンクを確認して移動してください。\" />")
			},
		},
		{
			name:      "Bot - Direct OGP when cushion disabled",
			shortUrl:  "bot-id",
			userAgent: "Twitterbot/1.0",
			setupMock: func(m *testutils.MockRedisClient) {
				m.On("GetBaseUrl", "bot-id").Return("https://example.com", nil).Once()
				m.On("GetIsNeedCusionPage", "bot-id").Return(false, nil).Once()
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, w *http.Response, body string) {
				// リダイレクト先のOGPを模した中間ページが表示されるはず
				assert.Contains(t, body, "https://example.com")
				assert.Contains(t, body, "<meta http-equiv=\"refresh\" content=\"0;url=https://example.com\">")
			},
		},
		{
			name:      "Bot - Cushion OGP when cushion enabled",
			shortUrl:  "bot-cushion-id",
			userAgent: "Slackbot 1.0",
			setupMock: func(m *testutils.MockRedisClient) {
				m.On("GetBaseUrl", "bot-cushion-id").Return("https://example.com", nil).Once()
				m.On("GetIsNeedCusionPage", "bot-cushion-id").Return(true, nil).Once()
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, w *http.Response, body string) {
				// クッションページのOGPが表示される
				assert.Contains(t, body, "<meta property=\"og:title\" content=\"リンクの確認 - URL Forge\" />")
			},
		},
		{
			name:     "Error - Not Found",
			shortUrl: "missing-id",
			setupMock: func(m *testutils.MockRedisClient) {
				m.On("GetBaseUrl", "missing-id").Return("", redis.Nil).Once()
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:     "Error - Redis GetBaseUrl failure",
			shortUrl: "error-id",
			setupMock: func(m *testutils.MockRedisClient) {
				m.On("GetBaseUrl", "error-id").Return("", errors.New("redis error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:     "Error - Redis GetIsNeedCusionPage failure",
			shortUrl: "error-cushion-id",
			setupMock: func(m *testutils.MockRedisClient) {
				m.On("GetBaseUrl", "error-cushion-id").Return("https://example.com", nil).Once()
				m.On("GetIsNeedCusionPage", "error-cushion-id").Return(false, errors.New("redis error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := new(testutils.MockRedisClient)
			tt.setupMock(mockRedis)

			appCtx := &utils.AppContext{
				Config: utils.Config{
					ServerEndpoint: "https://srv.test",
					AppName:        "URL Forge",
					BotUserAgents:  []string{"bot", "crawler", "spider", "facebookexternalhit", "twitterbot", "slackbot", "discordbot", "whatsapp", "line-poker"},
				},
				Redis: mockRedis,
			}

			// Gin のルーターにテンプレートをロードする必要がある（クッションページ用）
			router := gin.New()
			// テスト実行時のカレントディレクトリが app/controllers であることを想定
			router.LoadHTMLGlob("../templates/*.html")
			router.GET("/:shortUrl", controllers.GetUrlHandler(appCtx))

			var headers map[string]string
			if tt.userAgent != "" {
				headers = map[string]string{"User-Agent": tt.userAgent}
			}
			w := performRequest(router, "GET", "/"+tt.shortUrl, headers, nil)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, w.Result(), w.Body.String())
			}
			mockRedis.AssertExpectations(t)
		})
	}
}

func TestGetUrlHandler_Integration(t *testing.T) {
	_, mr, router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Run("Redirect Without Cushion", func(t *testing.T) {
		shortUrl := "no-cushion"
		baseUrl := "https://example.com/target"

		mr.HSet(shortUrl, "base_url", baseUrl)
		mr.HSet(shortUrl, "cushion", "false")

		w := performRequest(router, "GET", "/"+shortUrl, nil, nil)

		if w.Code != http.StatusFound {
			t.Errorf("expected status 302, got %d", w.Code)
		}
		if w.Header().Get("Location") != baseUrl {
			t.Errorf("expected redirect to %s, got %s", baseUrl, w.Header().Get("Location"))
		}
	})

	t.Run("Show Cushion Page", func(t *testing.T) {
		shortUrl := "with-cushion"
		baseUrl := "https://example.com/cushion-target"

		mr.HSet(shortUrl, "base_url", baseUrl)
		mr.HSet(shortUrl, "cushion", "true")

		w := performRequest(router, "GET", "/"+shortUrl, nil, nil)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), baseUrl) {
			t.Errorf("expected body to contain %s", baseUrl)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		w := performRequest(router, "GET", "/unknown", nil, nil)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}
