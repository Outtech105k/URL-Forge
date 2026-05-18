package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/Outtech105k/ShortUrlServer/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func ControlUrlHandler(appCtx *utils.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortUrlId := url.PathEscape(c.Param("shortUrl"))

		// Redisに問い合わせてURL情報を取得
		baseUrl, err := appCtx.Redis.GetBaseUrl(shortUrlId)
		if err != nil {
			if err == redis.Nil {
				c.HTML(http.StatusNotFound, "notfound.html", gin.H{
					"ServerEndpoint": appCtx.Config.ServerEndpoint,
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// 管理画面が公開設定になっているか確認
		isPublic, err := appCtx.Redis.GetIsPublicCtrl(shortUrlId)
		if err != nil {
			log.Printf("Failed to check public_ctrl: %v", err)
			// エラー時は安全のため非公開扱い
			isPublic = false
		}

		if !isPublic {
			c.HTML(http.StatusForbidden, "notfound.html", gin.H{
				"Message":        "このURLの管理画面は非公開です。",
				"ServerEndpoint": appCtx.Config.ServerEndpoint,
			})
			return
		}

		// OGP情報を取得
		ogp, err := utils.FetchOGPInfo(baseUrl)
		if err != nil {
			log.Printf("Failed to fetch OGP for %s: %v", baseUrl, err)
			// OGP取得失敗時は最低限の情報で表示
		}

		c.HTML(http.StatusOK, "control.html", gin.H{
			"URL":            baseUrl,
			"FullShortURL":   fmt.Sprintf("%s/%s", appCtx.Config.ServerEndpoint, shortUrlId),
			"OGPTitle":       ogp.Title,
			"OGPDescription": ogp.Description,
			"OGPImage":       ogp.Image,
			"Domain":         ogp.Domain,
		})
	}
}
