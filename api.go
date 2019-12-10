package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type tweetView struct {
	ID        int64    `json:"id"`
	User      string   `json:"user"`
	Text      string   `json:"text"`
	CreatedAt string   `json:"created_at"`
	Hashtags  []string `json:"hashtags"`
}

type stalkerAPI struct {
	stream  *tweetStream
	storage *tweetStorage
}

func newStalkerAPI(stream *tweetStream, storage *tweetStorage) (*stalkerAPI, error) {
	if stream == nil || storage == nil {
		return nil, errors.New("Stalker API can't be instanced without 'tweetStream' and 'tweetStorage'")
	}

	return &stalkerAPI{stream, storage}, nil
}

func (api *stalkerAPI) tweetsHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := u.Query()
	idParam := params.Get("id")
	if idParam == "" {
		return
	}
	id, _ := strconv.ParseInt(idParam, 10, 64)
	fmt.Println("ID : ", id)

	tweets := api.storage.QueryRecentByID(id)
	fmt.Println("Tweets: ", tweets)
	views := make([]tweetView, 0)
	hashtags := []string{"placeholder", "lorem", "ipsum"}
	for _, tweet := range tweets {
		tweetView := tweetView{
			tweet.ID,
			tweet.User.Name,
			tweet.FullText,
			tweet.CreatedAt,
			hashtags,
		}

		views = append(views, tweetView)
	}
	json.NewEncoder(w).Encode(views)
}

func (api *stalkerAPI) hashtagsHandler(w http.ResponseWriter, r *http.Request) {
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
	api.stream.Hashtags = hashtags

	fmt.Println("Setting new hashtags: ", api.stream.Hashtags)
}
