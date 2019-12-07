package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
)

func main() {
	hashtags := []string{"#MrRobot"}

	tweetChan := make(chan twitter.Tweet, 30)

	tweetStream := newTweetStream(newTwitterClientFromEnv(), hashtags, tweetChan)
	tweetStorage := newTweetStorage(hashtags, tweetChan)

	go tweetStream.Write()
	go tweetStorage.Read()

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	osChan := make(chan os.Signal)
	signal.Notify(osChan, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-osChan)
}
