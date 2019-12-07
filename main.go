package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func NewTwitterClient() *twitter.Client {
	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))

	if config.ConsumerKey == "" || config.ConsumerSecret == "" {
		log.Fatal("'TWITTER_CONSUMER_KEY' and 'TWITTER_CONSUMER_SECRET' env variables weren't set")
	} else if token.Token == "" || token.TokenSecret == "" {
		log.Fatal("'TWITTER_ACCESS_TOKEN' and 'TWITTER_ACCESS_TOKEN_SECRET' env variables weren't set")
	}

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	return client
}

type TweetStream struct {
	client          *twitter.Client
	hashtags        []string
	tweetChan       chan twitter.Tweet
	mostRecentTweet *twitter.Tweet
}

func NewTweetStream(client *twitter.Client, hashtags []string, outTweetChan chan twitter.Tweet) *TweetStream {
	return &TweetStream{client, hashtags, outTweetChan, nil}
}

func (tweetStream TweetStream) QueryRecentTweets(latestTweetID int64) ([]twitter.Tweet, int64) {
	query := strings.Join(tweetStream.hashtags[:], " OR ")
	searchParams := twitter.SearchTweetParams{
		Query:     query,
		Count:     30,
		TweetMode: "extended",
		SinceID:   latestTweetID,
	}

	search, _, err := tweetStream.client.Search.Tweets(&searchParams)

	if err != nil {
		log.Fatal(err)
	}

	return search.Statuses, search.Metadata.MaxID
}

func (tweetStream TweetStream) Write() {
	fmt.Println("Writing Stream...")

	latestTweetID := int64(0)

	for {
		tweets, maxID := tweetStream.QueryRecentTweets(latestTweetID)
		latestTweetID = maxID

		for _, tweet := range tweets {
			tweetStream.tweetChan <- tweet
		}
	}
}

type TweetStorage struct {
	hashtags  []string
	tweets    map[string][]*twitter.Tweet
	tweetChan chan twitter.Tweet
}

func NewTweetStorage(hashtags []string, inTweetChan chan twitter.Tweet) *TweetStorage {
	tweets := make(map[string][]*twitter.Tweet)

	return &TweetStorage{hashtags, tweets, inTweetChan}
}

func (storage TweetStorage) Read() {
	fmt.Println("Writeing Storage...")
	for {
		mainTweet := <-storage.tweetChan

		tweet := &mainTweet

		fmt.Println(tweet.User.Name)
		fmt.Println(tweet.Text)
		fmt.Println(tweet.CreatedAt)
		fmt.Println(tweet.Entities.Hashtags)
		fmt.Println("---------------------------------------------------------------------------")
		// storage.tweets
	}
}

func (storage TweetStorage) QueryTweets(hashtags []string, date time.Time) []*TweetStorage {
	for {
		// tweet := <-storage.tweetChan

		// fmt.Println(tweet.User.Name)
		// fmt.Println(tweet.Text)
		// fmt.Println(tweet.CreatedAt)

		// fmt.Println("---------------------------------------------------------------------------")
	}
}

func main() {
	hashtags := []string{"#MrRobot"}

	tweetChan := make(chan twitter.Tweet, 20)

	tweetStorage := NewTweetStorage(hashtags, tweetChan)
	go tweetStorage.Read()

	tweetStream := NewTweetStream(NewTwitterClient(), hashtags, tweetChan)
	go tweetStream.Write()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	osChan := make(chan os.Signal)
	signal.Notify(osChan, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-osChan)
}
