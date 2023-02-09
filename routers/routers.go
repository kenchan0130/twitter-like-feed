package routers

import (
	"fmt"
	"github.com/coocood/freecache"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/kenchan0130/twitter-like-feed/controllers"
	"github.com/kenchan0130/twitter-like-feed/repositories"
	"net/http"
	"os"
	"strconv"

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

	cacheExpiresSecStr := os.Getenv("CACHE_EXPIRES_SECONDS")
	var cacheExpiresSec int
	if cacheExpiresSecStr == "" {
		cacheExpiresSec = 60 * 60 * 2 // 60 * 60 * 2 is 2 hours
	} else {
		cacheExpiresSec, _ = strconv.Atoi(cacheExpiresSecStr)
	}

	rssRepository := repositories.RssRepository{
		Cache:         freecache.NewCache(100 * 1024 * 1024), // 100 * 1024 * 1024 is 100MB
		ExpireSeconds: cacheExpiresSec,
	}
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
		RssRepository:     rssRepository,
	}
	r.HEAD("/feed/:username", feedCtrl.Show)
	r.GET("/feed/:username", feedCtrl.Show)

	return r
}
