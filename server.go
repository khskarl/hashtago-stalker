package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

type apiServer struct {
	api  *stalkerAPI
	port string
	mux  *http.ServeMux
}

func newAPIServer(api *stalkerAPI) (*apiServer, error) {
	if api == nil {
		return nil, errors.New("API server can't be instanced without a 'stalker API'")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Listening at port ", port)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("assets/")))
	mux.HandleFunc("/tweets", api.tweetsHandler)
	mux.HandleFunc("/hashtags", api.hashtagsHandler)

	return &apiServer{api, port, mux}, nil
}

func (server *apiServer) listen() {
	http.ListenAndServe(":"+server.port, server.mux)
}
