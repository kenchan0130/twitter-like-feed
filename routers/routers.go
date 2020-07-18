package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/health")
	})
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.HEAD("/feed/:username", FeedUsernameGetHandler)
	r.GET("/feed/:username", FeedUsernameGetHandler)

	return r
}
