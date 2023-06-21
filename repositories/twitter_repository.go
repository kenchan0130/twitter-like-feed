package repositories

import (
	"context"
	"fmt"
	"github.com/kenchan0130/twitter-like-feed/models"
	"github.com/kenchan0130/unofficial-twitter-api-client-go/twitter"
	"github.com/samber/lo"
)

type TwitterRepository struct {
	Client twitter.Client
}

func (t TwitterRepository) GetLikesBy(screenName string) ([]models.Tweet, error) {
	userId, err := t.getUserIdBy(screenName)
	if err != nil {
		return nil, err
	}

	likeTweets, err := t.Client.GetUserLikingTweets(context.Background(), userId, 20)
	if err != nil {
		return nil, fmt.Errorf("user likes lookup error: %v", err)
	}

	tweets := lo.Map(likeTweets, func(tweet twitter.Tweet, _ int) models.Tweet {
		return models.Tweet{
			ID:               tweet.ID,
			Text:             tweet.Text,
			CreatedAt:        tweet.CreatedAt,
			AuthorID:         tweet.AuthorID,
			AuthorScreenName: tweet.AuthorScreenName,
		}
	})

	return tweets, nil
}

func (t TwitterRepository) getUserIdBy(screenName string) (string, error) {
	user, err := t.Client.GetUserByScreenName(context.Background(), screenName)
	if err != nil {
		return "", fmt.Errorf("user name lookup error: %v", err)
	}

	return user.ID, nil
}
