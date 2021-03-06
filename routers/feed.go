package routers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
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

var imgRe = regexp.MustCompile(`(?mi)<img(.*?)alt="(.+?)"(.*?)>`)
var ankerLinkRe = regexp.MustCompile(`(?mi)<a(.*?)href="(.+?)"(.*?)>(.*?)(</a>)`)
var hasktagRe = regexp.MustCompile(`(?mi)https://twitter.com/hashtag/(.+)\?\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

func parseTweetTextHTML(str string) string {
	extractImageStr := strings.Replace(strings.Replace(imgRe.ReplaceAllString(str, "$2"), "<br/>", "\n", -1), "<br>", "\n", -1)
	extractAnkerLinkStr := ankerLinkRe.ReplaceAllString(extractImageStr, "$2")
	extractHashtagStr := hasktagRe.ReplaceAllString(extractAnkerLinkStr, "#$1")
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(extractHashtagStr))
	return doc.Text()
}

func getTwitterLike(username string) (*[]Tweet, error) {
	url := fmt.Sprintf("https://syndication.twitter.com/timeline/likes?dnt=false&suppress_response_codes=true&screen_name=%s", username)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("An error occurred while trying to access %s, err: %s", url, err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("In accessing %s, HTTP status code is %d", url, res.StatusCode)
	}

	var twitterLikesResponse TwitterLikesResponse
	err = json.NewDecoder(res.Body).Decode(&twitterLikesResponse)
	if err != nil {
		return nil, fmt.Errorf("%s response returned an unexpected JSON. err: %s\n\n%s", url, err.Error(), res.Body)
	}

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(twitterLikesResponse.Body))

	tweetList := make([]Tweet, 0)
	tweetDateTimeLayout := "2006-01-02T15:04:05+0000"
	hasInvalidResponse := false
	doc.Find(".timeline-TweetList-tweet").Each(func(i int, s *goquery.Selection) {
		tweetAuthorName := strings.TrimSpace(s.Find(".TweetAuthor-name").Text())
		tweetAuthorScreenName := strings.TrimSpace(s.Find(".TweetAuthor-screenName").Text()) // With @ mark
		tweetTextHTML, err := s.Find(".timeline-Tweet-text").Html()

		if err != nil {
			log.Println(err)
			hasInvalidResponse = true
		}

		tweetText := parseTweetTextHTML(tweetTextHTML)
		log.Printf("TweetText: %s\n", tweetText)
		tweetURL := strings.TrimSpace(s.Find(".timeline-Tweet-timestamp").AttrOr("href", ""))
		tweetDateTime, err := time.Parse(tweetDateTimeLayout, strings.TrimSpace(s.Find(".dt-updated").AttrOr("datetime", "")))

		if err != nil {
			log.Println(err)
			hasInvalidResponse = true
		}

		if len(tweetAuthorName) == 0 || len(tweetAuthorScreenName) == 0 || len(tweetText) == 0 || len(tweetURL) == 0 {
			log.Printf("TweetAutorName: %s\n", tweetAuthorName)
			log.Printf("TweetAuthorScreenName: %s\n", tweetAuthorScreenName)
			log.Printf("TweetText: %s\n", tweetText)
			log.Printf("TweetURL: %s\n", tweetURL)
			hasInvalidResponse = true
		}

		tweet := Tweet{
			AuthorName:       tweetAuthorName,
			AuthorScreenName: tweetAuthorScreenName,
			Text:             tweetText,
			URL:              tweetURL,
			DateTime:         tweetDateTime,
		}
		tweetList = append(tweetList, tweet)
	})

	if hasInvalidResponse {
		return nil, fmt.Errorf("%s response returned an unexpected HTML in body attribute.\n\n%s", url, twitterLikesResponse.Body)
	}

	return &tweetList, nil
}

func generateFeed(username string, tweetList []Tweet) (string, error) {
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
	feed.Items = make([]*feeds.Item, 0)

	for _, tweet := range tweetList {
		item := &feeds.Item{
			Title:       fmt.Sprintf("@%s did LIKE a tweet of %s", username, tweet.AuthorScreenName),
			Link:        &feeds.Link{Href: tweet.URL},
			Description: fmt.Sprintf("@%s did LIKE %s tweet.", username, tweet.URL),
			Created:     tweet.DateTime,
			Id:          tweet.URL,
		}
		feed.Items = append(feed.Items, item)
	}

	return feed.ToRss()
}

// FeedUsernameGetHandler is function which returns Twitter like as feed
func FeedUsernameGetHandler(c *gin.Context) {
	username := strings.Replace(strings.TrimSpace(c.Param("username")), "@", "", 1)
	tweetList, err := getTwitterLike(username)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Error. Please check server log.")
		return
	}

	rss, err := generateFeed(username, *tweetList)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Error. Please check server log.")
		return
	}

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(http.StatusOK, rss)
}
