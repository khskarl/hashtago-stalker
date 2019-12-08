package main

import (
	"fmt"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

func printStorageInfo(storage *tweetStorage) {
	fmt.Println("***")
	for _, tweet := range storage.tweets {
		time, _ := tweet.CreatedAtTime()
		fmt.Println(time)
	}
	fmt.Println("***")
}

func main() {
	hashtags := []string{"#treebybike"}

	tweetStream := newTweetStream(newTwitterClientFromEnv(), hashtags)
	tweetStorage := newTweetStorage()

	storeChan := make(chan twitter.Tweet, 30)

	go tweetStream.Write(storeChan)
	go tweetStorage.Read(storeChan)

	for {
		go printStorageInfo(tweetStorage)
		var date time.Time
		for {
			var dateStr string
			fmt.Scanln(&dateStr)

			// str := 2019-12-08T05:31:41.000Z
			var err error
			date, err = time.Parse(time.RFC3339, dateStr)
			if err == nil {
				break
			}
		}
		tweets := tweetStorage.QueryByTime(date)
		fmt.Println(tweets)
	}

	// Wait for SIGINT and SIGTER5M (HIT CTRL-C)
	// osChan := make(chan os.Signal)
	// signal.Notify(osChan, syscall.SIGINT, syscall.SIGTERM)
	// log.Println(<-osChan)
}
