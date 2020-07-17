package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

type Tweet struct {
	AuthorName       string
	AuthorScreenName string
	Text             string
	URL              string
	DateTime         time.Time
}

type TwitterLikesResponse struct {
	Headers struct {
		Status      int32  `json:"status"`
		MaxPosition string `json:"maxPosition"`
		MinPosition string `json:"minPosition"`
		XPolling    int32  `json:"XPolling"`
		Time        int64  `json:"time"`
	} `json:"headers"`
	Body string `json:"body"`
}

func getTwitterLike(username string, lang string) (*[]Tweet, error) {
	res, err := http.Get(fmt.Sprintf("https://syndication.twitter.com/timeline/likes?lang=%s&screen_name=%s", lang, username))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP status code is %d", res.StatusCode)
	}

	var twitterLikesResponse TwitterLikesResponse
	err = json.NewDecoder(res.Body).Decode(&twitterLikesResponse)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(twitterLikesResponse.Body))
	if err != nil {
		return nil, fmt.Errorf("Unexpected HTML. err: %s\n\n%s", err.Error(), twitterLikesResponse.Body)
	}

	var tweetList = make([]Tweet, 0)
	doc.Find(".timeline-TweetList-tweet").Each(func(i int, s *goquery.Selection) {
		tweetAuthorName := s.Find(".TweetAuthor-name").Text()
		tweetAuthorScreenName := s.Find(".TweetAuthor-screenName").Text() // With @ mark
		tweetText := s.Find(".timeline-Tweet-text").Text()
		tweetURL := s.Find(".timeline-Tweet-timestamp").AttrOr("href", "")
		tweetDateTime, _ := time.Parse("2020-07-14T03:30:09+0000", s.Find(".dt-updated").AttrOr("datetime", ""))

		tweet := Tweet{
			AuthorName:       tweetAuthorName,
			AuthorScreenName: tweetAuthorScreenName,
			Text:             tweetText,
			URL:              tweetURL,
			DateTime:         tweetDateTime,
		}
		tweetList = append(tweetList, tweet)
	})

	return &tweetList, nil
}

func generateFeed(username string, tweetList []Tweet) (string, error) {
	feed := &feeds.Feed{
		Title:       fmt.Sprintf("@%s like feed | Twitter Like Feed", username),
		Link:        &feeds.Link{Href: fmt.Sprintf("https://twitter.com/%s/likes", username)},
		Description: fmt.Sprintf("@%s updated like list feed.", username),
	}
	feed.Items = make([]*feeds.Item, 0)

	for _, tweet := range tweetList {
		item := &feeds.Item{
			Title:       fmt.Sprintf("@%s did LIKE a tweet of %s.", username, tweet.AuthorScreenName),
			Link:        &feeds.Link{Href: fmt.Sprintf("https://twitter.com/%s/likes", username)},
			Description: fmt.Sprintf("@%s did LIKE %s tweet.", username, tweet.URL),
			Created:     tweet.DateTime,
			Id:          tweet.URL,
		}
		feed.Items = append(feed.Items, item)
	}

	return feed.ToRss()
}

func main() {
	port := 8080
	if len(os.Args) > 1 {
		p, _ := strconv.Atoi(os.Args[1])
		port = p
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/feed/:username", func(c *gin.Context) {
		username := c.Param("username")
		tweetList, err := getTwitterLike(username, "ja")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		rss, err := generateFeed(username, *tweetList)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.Header("Content-Type", "application/xml; charset=utf-8")
		c.String(http.StatusOK, rss)
	})
	r.Run(fmt.Sprintf(":%d", port))
}
