package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func newTwitterClientFromEnv() *twitter.Client {
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

// Helper function to smartly sleep coroutines according to the `RateLimitResource`
func sleepPerRateLimit(rateLimit *twitter.RateLimitResource) {
	if rateLimit.Remaining == 0 {
		timeToReset := time.Unix(int64(rateLimit.Reset), 0)
		durationUntilReset := time.Now().Sub(timeToReset)
		time.Sleep(durationUntilReset)
	} else {
		requestsPerMin := rateLimit.Limit / 15
		secsBetweenRequests := 60 / requestsPerMin
		sleepDuration := time.Duration(secsBetweenRequests * 1000000000)
		time.Sleep(sleepDuration)
	}
}

type tweetStream struct {
	client          *twitter.Client
	Hashtags        []string
	mostRecentTweet *twitter.Tweet
	rateLimits      *twitter.RateLimit
}

func newTweetStream(client *twitter.Client) *tweetStream {
	return &tweetStream{client, make([]string, 0), nil, nil}
}

func (stream *tweetStream) findAndSetMostRecentTweet(tweets []twitter.Tweet) {
	for _, tweet := range tweets {
		if tweet.ID > stream.mostRecentTweet.ID {
			stream.mostRecentTweet = &tweet
		}
	}
}

func (stream *tweetStream) fetchRecentTweets(mostRecentID int64) []twitter.Tweet {
	if len(stream.Hashtags[:]) == 0 {
		fmt.Println("[Stream] No hashtags defined, fetching no tweets")
		return make([]twitter.Tweet, 0)
	}

	hashtagsWithSymbol := make([]string, 0)
	for _, hashtagText := range stream.Hashtags {
		hashtagsWithSymbol = append(hashtagsWithSymbol, string('#')+hashtagText)
	}
	query := strings.Join(hashtagsWithSymbol[:], " OR ")
	// fmt.Println("[Stream] Fetching tweets. Since tweet ID: ", mostRecentID)

	// TODO: Find out if it's possible to exclude retweets
	// as we don't want retweets.
	searchParams := twitter.SearchTweetParams{
		Query:     query,
		Count:     10,
		TweetMode: "extended",
		SinceID:   mostRecentID,
	}
	search, _, err := stream.client.Search.Tweets(&searchParams)

	if err != nil {
		log.Fatal(err)
	}

	return search.Statuses
}

func (stream *tweetStream) UpdateRateLimits() {
	rateLimits, _, err := stream.client.RateLimits.Status(&twitter.RateLimitParams{Resources: []string{"application", "search"}})

	if err != nil {
		log.Fatal(err)
	}

	stream.rateLimits = rateLimits
}

func (stream *tweetStream) Write(outTweetChan chan twitter.Tweet) {
	rateLimitRoutine := func() {
		for {
			stream.UpdateRateLimits()

			rateStatusLimit := stream.rateLimits.Resources.Application["/application/rate_limit_status"]
			sleepPerRateLimit(rateStatusLimit)
		}
	}

	// Call it once so we have a `stream.rateLimits` instead of `nil`
	stream.UpdateRateLimits()
	// Coroutine for reading current rate limits
	go rateLimitRoutine()

	for {

		mostRecentID := int64(0)
		if stream.mostRecentTweet != nil {
			mostRecentID = stream.mostRecentTweet.ID
		}

		tweets := stream.fetchRecentTweets(mostRecentID)
		// fmt.Println("[Stream] Fetched: ", len(tweets), " tweets")

		if len(tweets) != 0 {
			if stream.mostRecentTweet == nil {
				stream.mostRecentTweet = &tweets[0]
			}

			stream.findAndSetMostRecentTweet(tweets)
		}

		for i := range tweets {
			// We `Write` the tweets in reverse because sending these tweets through the channel
			// ends up inverting the slice.
			// Originally `tweets` is ordered by decreasing date.
			tweet := tweets[len(tweets)-1-i]

			// Fetch the original tweet if it's a retweet
			if tweet.RetweetedStatus == nil {
				outTweetChan <- tweet

				// tweet = *tweet.RetweetedStatus
			}
		}

		tweetLimit := stream.rateLimits.Resources.Search["/search/tweets"]
		sleepPerRateLimit(tweetLimit)
	}
}

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

		if hasHashtag && tweet.ID > referenceID {
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
