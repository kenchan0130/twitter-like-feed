package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/kenchan0130/twitter-like-feed/models"
	"github.com/kenchan0130/twitter-like-feed/repositories"
	"github.com/samber/lo"
	"log"
	"net/http"
	"strings"
)

type FeedController struct {
	TwitterRepository repositories.TwitterRepository
	RssRepository     repositories.RssRepository
}

func (fc FeedController) Show(c *gin.Context) {
	username := strings.Replace(strings.TrimSpace(c.Param("username")), "@", "", 1)

	cachedRss := fc.RssRepository.GetBy(username)
	if cachedRss != nil {
		c.Header("Content-Type", "application/xml; charset=utf-8")
		c.String(http.StatusOK, *cachedRss)
		return
	}

	tweetList, err := fc.TwitterRepository.GetLikesBy(username)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Error. Please check server log.")
		return
	}

	rss, err := generateFeed(username, tweetList)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Error. Please check server log.")
		return
	}

	err = fc.RssRepository.SetBy(username, rss)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Error. Please check server log.")
		return
	}

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, rss)
}

func generateFeed(username string, tweetList []*models.Tweet) (string, error) {
	feed := &feeds.Feed{
		Title:       fmt.Sprintf("@%s like feed | Twitter Like Feed", username),
		Link:        &feeds.Link{Href: fmt.Sprintf("https://twitter.com/%s/likes", username)},
		Description: fmt.Sprintf("@%s updated like list feed.", username),
		Image: &feeds.Image{
			Link:  fmt.Sprintf("https://twitter.com/%s/likes", username),
			Url:   "https://abs.twimg.com/responsive-web/web/icon-default.3c3b2244.png", // From https://twitter.com/manifest.json
			Title: fmt.Sprintf("@%s like feed | Twitter Like Feed", username),
		},
	}

	feed.Items = lo.Map(tweetList, func(tweet *models.Tweet, _ int) *feeds.Item {
		tweetURL := fmt.Sprintf("https://twitter.com/%s/status/%s", tweet.AuthorScreenName, tweet.ID)
		return &feeds.Item{
			Title:       fmt.Sprintf("@%s did LIKE a tweet of %s", username, tweet.AuthorScreenName),
			Link:        &feeds.Link{Href: tweetURL},
			Description: fmt.Sprintf("@%s did LIKE %s tweet.", username, tweetURL),
			Created:     tweet.CreatedAt,
			Id:          strings.Join([]string{tweetURL, username}, "+"),
		}
	})

	return feed.ToRss()
}
