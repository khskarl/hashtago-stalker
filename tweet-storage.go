package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

type tweetStorage struct {
	tweets []twitter.Tweet
}

func newTweetStorage() *tweetStorage {
	tweets := make([]twitter.Tweet, 0)
	return &tweetStorage{tweets}
}

func (storage *tweetStorage) Read(inTweetChan chan twitter.Tweet) {
	for {
		tweet := <-inTweetChan

		storage.tweets = append(storage.tweets, tweet)
		sort.Slice(storage.tweets, func(i, j int) bool {
			return storage.tweets[i].ID < storage.tweets[j].ID
		})

		fmt.Println("[Storage] Reading new tweet. New storage size: ", len(storage.tweets))
		printTweet(&tweet)
		fmt.Println("---------------------------------------------------------------------------")
	}
}

func (storage tweetStorage) QuerySinceID(referenceID int64, hashtags []string) []twitter.Tweet {
	fmt.Println("[Storage] Querying since ID ", referenceID, " with hashtags: ", hashtags)
	filteredTweets := make([]twitter.Tweet, 0)
	for _, tweet := range storage.tweets {
		hasHashtag := false

		for _, hashtag := range hashtags {
			for _, hashtagEntity := range tweet.Entities.Hashtags {
				if strings.ToLower(hashtagEntity.Text) == strings.ToLower(hashtag) {
					hasHashtag = true
				}
			}
		}

		if hasHashtag && tweet.ID > referenceID && tweet.RetweetedStatus == nil {
			filteredTweets = append(filteredTweets, tweet)
		}
	}

	fmt.Println("[Storage] <Filtered Tweets>")
	for _, t := range filteredTweets {
		printTweet(&t)
	}

	return filteredTweets
}

func printTweet(tweet *twitter.Tweet) {
	hashtags := make([]string, 0)
	for _, hashtagEntity := range tweet.Entities.Hashtags {
		hashtags = append(hashtags, hashtagEntity.Text)
	}

	createdAt, _ := tweet.CreatedAtTime()
	offset, _ := time.ParseDuration("-03.00h")
	fmt.Println(tweet.User.Name)
	// fmt.Println(tweet.FullText)
	fmt.Println(tweet.IDStr)
	fmt.Println(hashtags)
	fmt.Println(createdAt.UTC().Add(offset))
	// fmt.Println(tweet.Entities.Hashtags)
}
