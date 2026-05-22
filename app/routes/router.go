package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/Outtech105k/ShortUrlServer/app/controllers"
	"github.com/Outtech105k/ShortUrlServer/app/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(appCtx *utils.AppContext) *gin.Engine {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	if appCtx.Config.AllowOrigins == "*" {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = strings.Split(appCtx.Config.AllowOrigins, ",")
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = false
	corsConfig.MaxAge = 12 * time.Hour
	r.Use(cors.New(corsConfig))

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"ServerEndpoint": appCtx.Config.ServerEndpoint,
			"AppName":        appCtx.Config.AppName,
		})
	})
	r.GET("/:shortUrl", controllers.GetUrlHandler(appCtx))
	r.GET("/:shortUrl/control", controllers.ControlUrlHandler(appCtx))
	r.POST("/api/set", controllers.SetUrlHandler(appCtx))

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "notfound.html", gin.H{
			"ServerEndpoint": appCtx.Config.ServerEndpoint,
			"AppName":        appCtx.Config.AppName,
		})
	})

	return r
}
