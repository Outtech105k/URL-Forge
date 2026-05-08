package controllers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/Outtech105k/ShortUrlServer/app/controllers"
	"github.com/Outtech105k/ShortUrlServer/app/models"
	"github.com/Outtech105k/ShortUrlServer/app/testutils"
	"github.com/Outtech105k/ShortUrlServer/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetUrlHandler_TableDriven(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// ヘルパー: 文字列へのポインタ作成
	ptrStr := func(s string) *string { return &s }
	ptrBool := func(b bool) *bool { return &b }
	ptrUint32 := func(u uint32) *uint32 { return &u }

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(m *testutils.MockRedisClient)
		expectedStatus int
		verifyResponse func(t *testing.T, body []byte)
	}{
		{
			name: "Success - Random ID",
			requestBody: models.SetUrlRequest{
				BaseURL: "https://example.com",
			},
			setupMock: func(m *testutils.MockRedisClient) {
				m.On("IsExists", mock.Anything).Return(false, nil).Once()
				m.On("SetURLRecord", mock.Anything, "https://example.com", false, mock.Anything).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, body []byte) {
				var resp models.APIResponse
				err := json.Unmarshal(body, &resp)
				assert.NoError(t, err)
				assert.Contains(t, resp.ShortURL, "https://srv.test/")
			},
		},
		{
			name: "Success - Custom ID",
			requestBody: models.SetUrlRequest{
				BaseURL:  "https://example.com",
				CustomID: ptrStr("my-custom-id"),
			},
			setupMock: func(m *testutils.MockRedisClient) {
				m.On("IsExists", "my-custom-id").Return(false, nil).Once()
				m.On("SetURLRecord", "my-custom-id", "https://example.com", false, mock.Anything).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Error - Custom ID Conflict",
			requestBody: models.SetUrlRequest{
				BaseURL:  "https://example.com",
				CustomID: ptrStr("existing-id"),
			},
			setupMock: func(m *testutils.MockRedisClient) {
				m.On("IsExists", "existing-id").Return(true, nil).Once()
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "Error - Parameter Conflict (CustomID with UseUppercase)",
			requestBody: models.SetUrlRequest{
				BaseURL:      "https://example.com",
				CustomID:     ptrStr("id"),
				UseUppercase: ptrBool(true),
			},
			setupMock:      func(m *testutils.MockRedisClient) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body []byte) {
				var apiErr models.APIError
				json.Unmarshal(body, &apiErr)
				assert.Equal(t, "parameter_conflict", apiErr.Type)
			},
		},
		{
			name: "Error - Invalid URL",
			requestBody: models.SetUrlRequest{
				BaseURL: "not-a-url",
			},
			setupMock:      func(m *testutils.MockRedisClient) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body []byte) {
				var apiErr models.APIError
				json.Unmarshal(body, &apiErr)
				assert.Equal(t, "validation_error", apiErr.Type)
			},
		},
		{
			name: "Success - Custom ID Length",
			requestBody: models.SetUrlRequest{
				BaseURL:  "https://example.com",
				IDLength: ptrUint32(8),
			},
			setupMock: func(m *testutils.MockRedisClient) {
				m.On("IsExists", mock.MatchedBy(func(id string) bool { return len(id) == 8 })).Return(false, nil).Once()
				m.On("SetURLRecord", mock.MatchedBy(func(id string) bool { return len(id) == 8 }), "https://example.com", false, mock.Anything).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Error - Custom ID contains slash",
			requestBody: models.SetUrlRequest{
				BaseURL:  "https://example.com",
				CustomID: ptrStr("invalid/id"),
			},
			setupMock:      func(m *testutils.MockRedisClient) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body []byte) {
				var apiErr models.APIError
				json.Unmarshal(body, &apiErr)
				assert.Contains(t, apiErr.Message, "contains `/`")
			},
		},
		{
			name: "Error - Redis IsExists failure",
			requestBody: models.SetUrlRequest{
				BaseURL:  "https://example.com",
				CustomID: ptrStr("id"),
			},
			setupMock: func(m *testutils.MockRedisClient) {
				m.On("IsExists", "id").Return(false, errors.New("redis error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Error - Redis SetURLRecord failure",
			requestBody: models.SetUrlRequest{
				BaseURL:  "https://example.com",
				CustomID: ptrStr("id"),
			},
			setupMock: func(m *testutils.MockRedisClient) {
				m.On("IsExists", "id").Return(false, nil).Once()
				m.On("SetURLRecord", "id", "https://example.com", false, mock.Anything).Return(errors.New("redis error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Error - ID generation failure (all 10 attempts exist)",
			requestBody: models.SetUrlRequest{
				BaseURL: "https://example.com",
			},
			setupMock: func(m *testutils.MockRedisClient) {
				m.On("IsExists", mock.Anything).Return(true, nil).Times(10)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Error - Malformed JSON",
			requestBody:    "invalid-json",
			setupMock:      func(m *testutils.MockRedisClient) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body []byte) {
				var apiErr models.APIError
				json.Unmarshal(body, &apiErr)
				assert.Equal(t, "invalid_request", apiErr.Type)
			},
		},
		{
			name:           "Error - Empty Body",
			requestBody:    nil,
			setupMock:      func(m *testutils.MockRedisClient) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body []byte) {
				var apiErr models.APIError
				json.Unmarshal(body, &apiErr)
				assert.Equal(t, "invalid_request", apiErr.Type)
				assert.Equal(t, "Empty JSON body", apiErr.Message)
			},
		},
		{
			name: "Error - No character types available",
			requestBody: models.SetUrlRequest{
				BaseURL:      "https://example.com",
				UseUppercase: ptrBool(false),
				UseLowercase: ptrBool(false),
				UseNumbers:   ptrBool(false),
			},
			setupMock:      func(m *testutils.MockRedisClient) {},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: func(t *testing.T, body []byte) {
				var apiErr models.APIError
				json.Unmarshal(body, &apiErr)
				assert.Equal(t, "invalid_request", apiErr.Type)
				assert.Equal(t, "No character types available for URL ID.", apiErr.Message)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := new(testutils.MockRedisClient)
			tt.setupMock(mockRedis)

			appCtx := &utils.AppContext{
				Config: utils.Config{
					ServerEndpoint: "https://srv.test",
				},
				Redis: mockRedis,
			}

			router := gin.New()
			router.POST("/set", controllers.SetUrlHandler(appCtx))

			w := performRequest(router, "POST", "/set", tt.requestBody)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, w.Body.Bytes())
			}
			mockRedis.AssertExpectations(t)
		})
	}
}

func TestSetUrlHandler_Integration(t *testing.T) {
	appCtx, mr, router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Run("Create Random URL", func(t *testing.T) {
		reqBody := models.SetUrlRequest{
			BaseURL: "https://example.com/original",
		}
		w := performRequest(router, "POST", "/set", reqBody)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp models.APIResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("responce json unmarshal: %v", err)
		}

		short := strings.TrimPrefix(resp.ShortURL, appCtx.Config.ServerEndpoint+"/")
		mrGot := mr.HGet(short, "base_url")
		if mrGot != resp.BaseURL {
			t.Errorf("expected base_url %s, got %s", resp.BaseURL, mrGot)
		}
	})

	t.Run("Success with Custom ID", func(t *testing.T) {
		customId := "my-id"
		reqBody := models.SetUrlRequest{
			BaseURL:  "https://example.com/custom",
			CustomID: &customId,
		}
		w := performRequest(router, "POST", "/set", reqBody)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		if !mr.Exists(customId) {
			t.Errorf("expected key %s to exist in redis", customId)
		}
		mrGot := mr.HGet(customId, "base_url")
		if mrGot != reqBody.BaseURL {
			t.Errorf("expected base_url %s, got %s", reqBody.BaseURL, mrGot)
		}
	})

	t.Run("Conflict with Existing ID", func(t *testing.T) {
		existingId := "already-used"
		mr.HSet(existingId, "base_url", "https://existing.com")

		reqBody := models.SetUrlRequest{
			BaseURL:  "https://example.com/new",
			CustomID: &existingId,
		}
		w := performRequest(router, "POST", "/set", reqBody)

		if w.Code != http.StatusConflict {
			t.Errorf("expected status 409, got %d", w.Code)
		}
	})

	t.Run("Validation Error (Invalid URL)", func(t *testing.T) {
		reqBody := models.SetUrlRequest{
			BaseURL: "not-a-url",
		}
		w := performRequest(router, "POST", "/set", reqBody)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}
