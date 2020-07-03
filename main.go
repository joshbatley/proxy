package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/joshbatley/proxy/client"
	"github.com/joshbatley/proxy/config"
	"github.com/joshbatley/proxy/query"
)

func main() {
	config, err := config.Load("./config.yml")

	if err != nil {
		panic("Config unreadable")
	}

	// DB setup

	r := mux.NewRouter()
	r.HandleFunc("/config", client.Serve)
	r.HandleFunc("/query", query.Serve)

	http.Handle("/", r)

	log.Println("Listing on localhosts:" + config.Port)
	http.ListenAndServe("localhost:"+config.Port, nil)
}
