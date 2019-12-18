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
	ID        string   `json:"id"`
	User      string   `json:"user"`
	Text      string   `json:"text"`
	CreatedAt int64    `json:"created_at"`
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

func (api *stalkerAPI) tweetsHandler(writer http.ResponseWriter, request *http.Request) {
	u, err := url.Parse(request.URL.String())
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	params := u.Query()
	idParam := params.Get("id")
	if idParam == "" {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	id, _ := strconv.ParseInt(idParam, 10, 64)
	fmt.Println("[Stalker Api] Querying tweets since ID: ", id)

	tweets := api.storage.QuerySinceID(id, api.stream.Hashtags)
	views := make([]tweetView, 0)
	for _, tweet := range tweets {
		time, _ := tweet.CreatedAtTime()
		hashtags := make([]string, 0)
		for _, hashtagEntity := range tweet.Entities.Hashtags {
			hashtags = append(hashtags, hashtagEntity.Text)
		}

		tweetView := tweetView{
			tweet.IDStr,
			tweet.User.Name,
			tweet.FullText,
			time.Unix(),
			hashtags,
		}

		views = append(views, tweetView)
	}

	json.NewEncoder(writer).Encode(views)
}

func (api *stalkerAPI) hashtagsHandler(writer http.ResponseWriter, request *http.Request) {
	type hashtagsView struct {
		Hashtags []string `json:"hashtags"`
	}

	var view hashtagsView
	err := json.NewDecoder(request.Body).Decode(&view)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	api.stream.Hashtags = view.Hashtags
	fmt.Println("[Stalker Api] Setting new hashtags: ", api.stream.Hashtags)
}
