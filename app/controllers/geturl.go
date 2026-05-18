package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

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

		// Bot判定（OGP取得用）
		userAgent := c.GetHeader("User-Agent")
		isBot := isBotUserAgent(userAgent)

		if isCushionRequired {
			// クッションページを表示（ネタバレ防止OGP）
			c.HTML(http.StatusOK, "cushion.html", gin.H{
				"URL":            baseUrl,
				"FullShortURL":   fmt.Sprintf("%s/%s", appCtx.Config.ServerEndpoint, shortUrl),
				"ServerEndpoint": appCtx.Config.ServerEndpoint,
			})
			return
		}

		if isBot {
			// クッションページなしだがBotの場合：リダイレクト先のOGPを返す
			ogp, err := utils.FetchOGPInfo(baseUrl)
			if err != nil {
				log.Printf("Failed to fetch OGP for bot redirect: %v", err)
				// 失敗時は最低限の情報で返す
			}

			c.HTML(http.StatusOK, "direct_ogp.html", gin.H{
				"URL":            baseUrl,
				"FullShortURL":   fmt.Sprintf("%s/%s", appCtx.Config.ServerEndpoint, shortUrl),
				"OGPTitle":       ogp.Title,
				"OGPDescription": ogp.Description,
				"OGPImage":       ogp.Image,
			})
			return
		}

		// クッションページなしでリダイレクト
		c.Redirect(http.StatusFound, baseUrl)
	}
}

func isBotUserAgent(ua string) bool {
	ua = strings.ToLower(ua)
	bots := []string{
		"bot",
		"crawler",
		"spider",
		"facebookexternalhit",
		"twitterbot",
		"slackbot",
		"discordbot",
		"whatsapp",
		"line-poker",
	}
	for _, bot := range bots {
		if strings.Contains(ua, bot) {
			return true
		}
	}
	return false
}
