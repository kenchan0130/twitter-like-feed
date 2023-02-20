package repositories

import (
	"context"
	"fmt"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/kenchan0130/twitter-like-feed/models"
	"github.com/samber/lo"
	"strings"
	"time"
)

type TwitterRepository struct {
	Client twitter.Client
}

func (t TwitterRepository) GetLikesBy(screenName string) ([]*models.Tweet, error) {
	userId, err := t.getUserIdBy(screenName)
	if err != nil {
		return nil, err
	}

	likeResponse, err := t.Client.UserLikesLookup(context.Background(), userId, twitter.UserLikesLookupOpts{
		Expansions:  []twitter.Expansion{twitter.ExpansionAuthorID},
		TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt},
	})
	if err != nil {
		return nil, fmt.Errorf("user likes lookup error: %v", err)
	}

	authorIdList := lo.Map(likeResponse.Raw.Tweets, func(tweetObj *twitter.TweetObj, _ int) string {
		return tweetObj.AuthorID
	})

	authorMapping, err := t.getUsersBy(authorIdList)
	if err != nil {
		return nil, err
	}

	tweets := lo.Map(likeResponse.Raw.Tweets, func(tweetObj *twitter.TweetObj, _ int) *models.Tweet {
		tweetDateTime, _ := time.Parse(time.RFC3339, tweetObj.CreatedAt)
		return &models.Tweet{
			ID:               tweetObj.ID,
			Text:             tweetObj.Text,
			CreatedAt:        tweetDateTime,
			AuthorID:         tweetObj.AuthorID,
			AuthorScreenName: authorMapping[tweetObj.AuthorID].User.UserName,
		}
	})

	return tweets, nil
}

func (t TwitterRepository) getUserIdBy(screenName string) (string, error) {
	userResponse, err := t.Client.UserNameLookup(context.Background(), []string{screenName}, twitter.UserLookupOpts{})
	if err != nil {
		return "", fmt.Errorf("user name lookup error: %v", err)
	}

	userDic := userResponse.Raw.UserDictionaries()
	if userDic == nil {
		return "", fmt.Errorf("user lookup error with %s", screenName)
	}
	keys := lo.Keys[string, *twitter.UserDictionary](userDic)
	id := keys[0]

	return id, nil
}

func (t TwitterRepository) getUsersBy(userIds []string) (map[string]*twitter.UserDictionary, error) {
	userResponse, err := t.Client.UserLookup(context.Background(), userIds, twitter.UserLookupOpts{
		Expansions: []twitter.Expansion{twitter.ExpansionPinnedTweetID},
	})
	if err != nil {
		return nil, fmt.Errorf("user lookup error: %v", err)
	}

	userDic := userResponse.Raw.UserDictionaries()
	if userDic == nil {
		return nil, fmt.Errorf("user lookup error with %s", strings.Join(userIds, ", "))
	}

	return userDic, nil
}
