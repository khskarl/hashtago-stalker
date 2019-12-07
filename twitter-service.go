package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func newTwitterClient() *twitter.Client {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("TWITTER_CONSUMER_KEY"),
		ClientSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}

	if config.ClientID == "" || config.ClientSecret == "" {
		log.Fatal("'TWITTER_CONSUMER_KEY' and 'TWITTER_CONSUMER_SECRET' env variables weren't set")
	}

	httpClient := config.Client(oauth2.NoContext)
	client := twitter.NewClient(httpClient)

	return client
}

type tweetStream struct {
	client          *twitter.Client
	hashtags        []string
	tweetChan       chan twitter.Tweet
	mostRecentTweet *twitter.Tweet
	rateLimits      *twitter.RateLimit
}

func newTweetStream(client *twitter.Client, hashtags []string, outTweetChan chan twitter.Tweet) *tweetStream {
	return &tweetStream{client, hashtags, outTweetChan, nil, nil}
}

func (tweetStream tweetStream) queryRecentTweets(latestTweetID int64) ([]twitter.Tweet, int64) {
	query := strings.Join(tweetStream.hashtags[:], " OR ")
	searchParams := twitter.SearchTweetParams{
		Query:     query,
		Count:     10,
		TweetMode: "extended",
		SinceID:   latestTweetID,
	}

	search, _, err := tweetStream.client.Search.Tweets(&searchParams)

	if err != nil {
		log.Fatal(err)
	}

	return search.Statuses, search.Metadata.MaxID
}

func (tweetStream tweetStream) Write() {
	go func() {
		for {
			fmt.Println("**QUERYING APP**")
			rateLimits, _, err := tweetStream.client.RateLimits.Status(&twitter.RateLimitParams{Resources: []string{"application", "search"}})

			if err != nil {
				log.Fatal(err)
			}

			tweetStream.rateLimits = rateLimits
			rateStatusLimit := tweetStream.rateLimits.Resources.Application["/application/rate_limit_status"]

			if rateStatusLimit.Remaining == 0 {
				timeToReset := time.Unix(int64(rateStatusLimit.Reset), 0)
				durationUntilReset := time.Now().Sub(timeToReset)
				fmt.Println("APP REM: ", durationUntilReset)
				time.Sleep(durationUntilReset)
			} else {
				requestsPerMin := rateStatusLimit.Limit / 15
				secsBetweenRequests := 60 / requestsPerMin
				sleepDuration := time.Duration(secsBetweenRequests * 1000000000)
				fmt.Println("APP: ", sleepDuration)
				time.Sleep(sleepDuration)
			}
		}
	}()

	latestTweetID := int64(0)

	for {
		fmt.Println("**QUERYING TWEETS**")

		tweets, maxID := tweetStream.queryRecentTweets(latestTweetID)
		latestTweetID = maxID

		for i := range tweets {
			tweet := tweets[len(tweets)-1-i]
			tweetStream.tweetChan <- tweet
		}

		tweetLimit := tweetStream.rateLimits.Resources.Search["/search/tweets"]
		if tweetLimit.Remaining == 0 {
			timeToReset := time.Unix(int64(tweetLimit.Reset), 0)
			durationUntilReset := time.Now().Sub(timeToReset)
			fmt.Println("TWEET REM: ", durationUntilReset)
			time.Sleep(durationUntilReset)
		} else {
			requestsPerMin := tweetLimit.Limit / 15
			secsBetweenRequests := 60 / requestsPerMin
			sleepDuration := time.Duration(secsBetweenRequests * 1000000000)
			fmt.Println("TWEET: ", sleepDuration)
			time.Sleep(sleepDuration)
		}
	}
}

type tweetStorage struct {
	hashtags  []string
	tweets    map[string][]*twitter.Tweet
	tweetChan chan twitter.Tweet
}

func newTweetStorage(hashtags []string, inTweetChan chan twitter.Tweet) *tweetStorage {
	tweets := make(map[string][]*twitter.Tweet)
	return &tweetStorage{hashtags, tweets, inTweetChan}
}

func (storage tweetStorage) Read() {
	for {
		mainTweet := <-storage.tweetChan

		tweet := &mainTweet

		fmt.Println(tweet.User.Name)
		fmt.Println(tweet.FullText)
		fmt.Println(tweet.CreatedAt)
		fmt.Println(tweet.Entities.Hashtags)
		fmt.Println("---------------------------------------------------------------------------")
	}
}
