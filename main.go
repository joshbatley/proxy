package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/config"
	"github.com/joshbatley/proxy/database"
	"github.com/joshbatley/proxy/handler"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config, err := config.Load("./config.yml")

	if err != nil {
		panic("Config unreadable")
	}

	// DB setup
	database.Conn()

	r := mux.NewRouter()
	r.HandleFunc("/config", handler.ClientServe)
	r.HandleFunc("/query", handler.QueryServe)
	http.Handle("/", r)

	log.Println("Listing on localhosts:" + config.Port)
	http.ListenAndServe("localhost:"+config.Port, nil)
}
