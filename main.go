package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/dghubble/go-twitter/twitter"
)

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

	fmt.Println("Hashtags are is: ", hashtags)
	// time := time.Now().AddDate(0, 0, -1)
	// tweets := tweetStorage.QueryByTime(time)
	json.NewEncoder(w).Encode("hey")
}

func main() {
	hashtags := []string{"#treebybike"}

	tweetStream := newTweetStream(newTwitterClientFromEnv(), hashtags)
	tweetStorage := newTweetStorage()

	storeChan := make(chan twitter.Tweet, 30)

	go tweetStream.Write(storeChan)
	go tweetStorage.Read(storeChan)

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
