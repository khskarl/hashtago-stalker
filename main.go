package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

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

	id, _ := strconv.ParseInt(params.Get("id"), 10, 64)
	fmt.Println("ID : ", id)

	tweets := storage.QueryRecentByID(id)
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

func hashtagsHandler(w http.ResponseWriter, r *http.Request) {
	type hashtagsView struct {
		Hashtags []string `json:"hashtags"`
	}

	var view hashtagsView
	err := json.NewDecoder(r.Body).Decode(&view)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashtags := view.Hashtags
	stream.Hashtags = hashtags

	fmt.Println("Setting new hashtags: ", stream.Hashtags)
}

var stream *tweetStream
var storage *tweetStorage

func main() {
	hashtags := []string{"#treebybike"}

	stream = newTweetStream(newTwitterClientFromEnv(), hashtags)
	storage = newTweetStorage()

	storeChan := make(chan twitter.Tweet, 30)

	go stream.Write(storeChan)
	go storage.Read(storeChan)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Listening at port ", port)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("assets/")))
	mux.HandleFunc("/tweets", tweetsHandler)
	mux.HandleFunc("/hashtags", hashtagsHandler)
	http.ListenAndServe(":"+port, mux)
}
