package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
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
	stream    *twitter.Stream
	demux     twitter.SwitchDemux
	tweetChan chan *twitter.Tweet
}

func NewTweetStream(client *twitter.Client, hashtags []string, outTweetChan chan *twitter.Tweet) *TweetStream {

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		outTweetChan <- tweet
	}

	filterParams := &twitter.StreamFilterParams{
		Track:         hashtags,
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	return &TweetStream{stream, demux, outTweetChan}
}

func (tweetStream TweetStream) Start() {
	fmt.Println("Starting Stream...")

	go tweetStream.demux.HandleChan(tweetStream.stream.Messages)
}

func (tweetStream TweetStream) Stop() {
	fmt.Println("Stopping Stream...")
	tweetStream.stream.Stop()
}

type TweetStorage struct {
	hashtags  []string
	tweets    map[string][]*twitter.Tweet
	tweetChan chan *twitter.Tweet
}

func NewTweetStorage(hashtags []string, inTweetChan chan *twitter.Tweet) *TweetStorage {
	tweets := make(map[string][]*twitter.Tweet)

	return &TweetStorage{hashtags, tweets, inTweetChan}
}

func (storage TweetStorage) ReadStream() {
	fmt.Println("Starting Storage...")
	for {
		mainTweet := <-storage.tweetChan

		tweet := mainTweet
		if mainTweet.Retweeted == true {
			tweet = mainTweet.RetweetedStatus
			fmt.Println("IM RETWEET")
		}
		extendedTweet := tweet.ExtendedTweet
		if extendedTweet != nil {
			fmt.Println("EXTENDED")
		}

		fmt.Println(tweet.User.Name)
		fmt.Println(tweet.Text)
		fmt.Println(tweet.CreatedAt)
		fmt.Println(tweet.Entities.Hashtags)
		if extendedTweet != nil {
			fmt.Println(extendedTweet.Entities.Hashtags)
		}
		if mainTweet.Retweeted == true {
			fmt.Println(mainTweet.Entities.Hashtags)
		}
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
	client := NewTwitterClient()
	hashtags := []string{"#MrRobot", "#sverige", "#svenska", "#brexit"}

	tweetChan := make(chan *twitter.Tweet, 20)
	tweetStorage := NewTweetStorage(hashtags, tweetChan)
	go tweetStorage.ReadStream()
	tweetStream := NewTweetStream(client, hashtags, tweetChan)

	tweetStream.Start()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	osChan := make(chan os.Signal)
	signal.Notify(osChan, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-osChan)

	tweetStream.Stop()

}
