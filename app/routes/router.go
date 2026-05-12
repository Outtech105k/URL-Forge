package routes

import (
	"net/http"
	"time"

	"github.com/Outtech105k/ShortUrlServer/app/controllers"
	"github.com/Outtech105k/ShortUrlServer/app/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(appCtx *utils.AppContext) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/:shortUrl", controllers.GetUrlHandler(appCtx))
	r.POST("/api/set", controllers.SetUrlHandler(appCtx))

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "notfound.html", nil)
	})

	return r
}
