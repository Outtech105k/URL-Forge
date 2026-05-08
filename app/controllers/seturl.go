package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/Outtech105k/ShortUrlServer/app/models"
	"github.com/Outtech105k/ShortUrlServer/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type NotAcceptableIdError struct {
	Message string
}

func (naie *NotAcceptableIdError) Error() string {
	return naie.Message
}

func SetUrlHandler(appCtx *utils.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r models.SetUrlRequest

		if err := c.ShouldBindJSON(&r); err != nil {

			// エラー判定 (JSON varidator のみで判断できる内容)
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				details := make([]map[string]string, 0)
				for _, fe := range ve {
					jsonField := extractJSONFieldName(r, fe.StructField())
					details = append(details, map[string]string{
						"field":   jsonField,
						"message": utils.ValidationErrorMessage(jsonField, fe.Tag()),
					})
				}

				apiErr := models.APIError{
					Type:    "validation_error",
					Details: details,
				}
				c.JSON(http.StatusBadRequest, apiErr)
				return
			}

			// JSONが空
			if errors.Is(err, io.EOF) {
				apiErr := models.APIError{
					Type:    "invalid_request",
					Message: "Empty JSON body",
				}
				c.JSON(http.StatusBadRequest, apiErr)
				return
			}

			// JSONの構文エラー
			var syntaxErr *json.SyntaxError
			if errors.As(err, &syntaxErr) {
				apiErr := models.APIError{
					Type:    "invalid_request",
					Message: "Malformed JSON body",
				}
				c.JSON(http.StatusBadRequest, apiErr)
				return
			}

			// その他のJSONバインドエラー
			apiErr := models.APIError{
				Type:    "invalid_request",
				Message: "Invalid JSON input",
			}
			c.JSON(http.StatusBadRequest, apiErr)
			return
		}

		if apiErr := setUrlHandlerCustomValidate(&r); apiErr != nil {
			c.JSON(http.StatusBadRequest, apiErr)
			return
		}

		// UseUppercase, UseLowercase, UseNumbers, IDLength, SandCushionのデフォルト値を設定
		// ExpireInはnilの場合、無期限として扱うのでnilを許す
		nilSetDefault(&r.UseUppercase, false)
		nilSetDefault(&r.UseLowercase, true)
		nilSetDefault(&r.UseNumbers, true)
		nilSetDefault(&r.IDLength, 6)
		nilSetDefault(&r.SandCushion, false)

		var customId string
		if r.CustomID == nil {
			// カスタムIDが指定されていない場合、4文字カスタムIDの生成（最大10回試行）
			customIdIsExists := false
			for i := 0; i < 10; i++ {
				var err error
				customId, err = utils.MakeRandomStr(
					*r.IDLength,
					*r.UseUppercase,
					*r.UseLowercase,
					*r.UseNumbers,
				)
				if err != nil {
					// ランダム生成に必要な文字がない場合
					if err == utils.ErrNoCharacterSet {
						c.JSON(http.StatusBadRequest, models.APIError{
							Type:    "invalid_request",
							Message: "No character types available for URL ID.",
						})
						return
					}

					// それ以外のランダム生成エラー
					returnInternalServerError(c, fmt.Sprintf("MakeRandomStr error: %v", err))
					return
				}

				// 受理可能なカスタムIDか
				if err := checkAcceptableUrlId(customId); err != nil {
					continue
				}

				// 生成されたカスタムIDがRedisに存在するか確認
				customIdIsExists, err = appCtx.Redis.IsExists(customId)
				if err != nil {
					returnInternalServerError(c, fmt.Sprintf("Redis generated ID exists error: %v", err))
					return
				}

				if !customIdIsExists {
					break
				}
			}
			if customIdIsExists {
				returnInternalServerError(c, "Custom ID generation failed after 10 attempts.")
				return
			}
		} else { // カスタムIDが指定されている場合
			// 受理可能なカスタムIDか
			if err := checkAcceptableUrlId(*r.CustomID); err != nil {
				c.JSON(http.StatusBadRequest, models.APIError{
					Type:    "invalid_request",
					Message: fmt.Sprintf("custom_id: %s", err.Error()),
				})

				return
			}

			// 適切にエスケープする
			customId = url.PathEscape(*r.CustomID)

			// Redisに存在するか確認
			customIdIsExists, err := appCtx.Redis.IsExists(customId)
			if err != nil {
				returnInternalServerError(c, fmt.Sprintf("Redis custom ID exists error: %v", err))
				return
			}

			if customIdIsExists {
				c.JSON(http.StatusConflict, models.APIError{
					Type:    "conflict",
					Message: "custome_id is already used.",
				})
				return
			}
		}

		// URLの有効期限を設定
		var expireIn *time.Duration = nil
		if r.ExpireIn != nil {
			expireIn = &r.ExpireIn.Duration
		}

		// RedisにURLを保存
		if err := appCtx.Redis.SetURLRecord(customId, r.BaseURL, *r.SandCushion, expireIn); err != nil {
			returnInternalServerError(c, fmt.Sprintf("Redis set URL record error: %v", err))
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			BaseURL:  r.BaseURL,
			ShortURL: fmt.Sprintf("%s/%s", appCtx.Config.ServerEndpoint, customId),
		})
	}
}

func setUrlHandlerCustomValidate(r *models.SetUrlRequest) *models.APIError {
	if r.CustomID != nil && (r.UseUppercase != nil || r.UseLowercase != nil || r.UseNumbers != nil || r.IDLength != nil) {
		return &models.APIError{
			Type:    "parameter_conflict",
			Message: "custom_id cannot be used together with use_uppercase, use_lowercase, use_numbers, or id_length",
		}
	}

	return nil
}

func checkAcceptableUrlId(s string) error {
	// `/` はGinの仕様上リダイレクト時に処理されないので禁止する
	if strings.Contains(s, "/") {
		return &NotAcceptableIdError{
			Message: "contains `/`, it isn't acceptable.",
		}
	}

	return nil
}

func nilSetDefault[T any](v **T, defaultV T) {
	if *v == nil {
		*v = &defaultV
	}
}

func returnInternalServerError(c *gin.Context, logMsg string) {
	c.JSON(http.StatusInternalServerError, models.APIError{
		Type:    "internal_error",
		Message: "An unexpected error occurred. Please try again later.",
	})
	log.Println(logMsg)
}

// JSONフィールド名を出力
func extractJSONFieldName(obj any, fieldName string) string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if f, ok := t.FieldByName(fieldName); ok {
		tag := f.Tag.Get("json")
		if tag != "" && tag != "-" {
			// json tag may contain ",omitempty", so split
			return strings.Split(tag, ",")[0]
		}
	}
	return fieldName // fallback
}
