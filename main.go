package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

type tweetView struct {
	ID        int64  `json:"id"`
	User      string `json:"user"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

func tweetsHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := u.Query()
	fmt.Println(params)
	hashtags := params.Get("q")
	fmt.Println("Hashtags are : ", hashtags)

	time := time.Now().AddDate(0, 0, -1)
	tweets := storage.QueryByTime(time)
	views := make([]tweetView, 0)
	for _, tweet := range tweets {
		tweetView := tweetView{
			tweet.ID,
			tweet.User.Name,
			tweet.FullText,
			tweet.CreatedAt,
		}

		views = append(views, tweetView)
	}
	json.NewEncoder(w).Encode(views)
}

var storage *tweetStorage

func main() {
	hashtags := []string{"#treebybike"}

	tweetStream := newTweetStream(newTwitterClientFromEnv(), hashtags)
	storage = newTweetStorage()

	storeChan := make(chan twitter.Tweet, 30)

	go tweetStream.Write(storeChan)
	go storage.Read(storeChan)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Listening at port ", port)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("assets/")))
	mux.HandleFunc("/tweets", tweetsHandler)
	http.ListenAndServe(":"+port, mux)
}
