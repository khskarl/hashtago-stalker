package main

import (
	"github.com/dghubble/go-twitter/twitter"
)

func main() {
	stream := newTweetStream(newTwitterClientFromEnv())
	stream.Hashtags = []string{}
	storage := newTweetStorage()

	storeChan := make(chan twitter.Tweet, 30)

	go stream.Write(storeChan)
	go storage.Read(storeChan)

	api, _ := newStalkerAPI(stream, storage)
	server, _ := newAPIServer(api)

	server.listen()
}
