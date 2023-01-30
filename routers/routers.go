package routers

import (
	"fmt"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/kenchan0130/twitter-like-feed/controllers"
	"github.com/kenchan0130/twitter-like-feed/repositories"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func Init() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/health")
	})
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	token := os.Getenv("BEARER_TOKEN")
	twitterRepository := repositories.TwitterRepository{
		Client: twitter.Client{
			Authorizer: authorize{
				Token: token,
			},
			Client: http.DefaultClient,
			Host:   "https://api.twitter.com",
		},
	}
	feedCtrl := controllers.FeedController{
		TwitterRepository: twitterRepository,
	}
	r.HEAD("/feed/:username", feedCtrl.Show)
	r.GET("/feed/:username", feedCtrl.Show)

	return r
}
