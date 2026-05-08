package controllers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/Outtech105k/ShortUrlServer/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func GetUrlHandler(appCtx *utils.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortUrl := url.PathEscape(c.Param("shortUrl"))

		// Redisに問い合わせてURLを取得
		baseUrl, err := appCtx.Redis.GetBaseUrl(shortUrl)
		if err != nil {
			// 保存されていない(nil)場合は404を返す
			if err == redis.Nil {
				c.HTML(http.StatusNotFound, "notfound.html", nil)

				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve base URL",
			})
			log.Printf("Failed to retrieve base URL: %v", err)

			return
		}

		// クッションページが必要か確認
		isCushionRequired, err := appCtx.Redis.GetIsNeedCusionPage(shortUrl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			log.Printf("Failed to check if cushion page is needed: %v", err)

			return
		}

		if isCushionRequired {
			// クッションページを表示
			c.HTML(http.StatusOK, "cushion.html", gin.H{
				"URL": baseUrl,
			})

			return
		}

		// クッションページなしでリダイレクト
		c.Redirect(http.StatusFound, baseUrl)
	}
}
