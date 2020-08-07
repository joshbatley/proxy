package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/database"
	"github.com/joshbatley/proxy/handler"
	"github.com/joshbatley/proxy/repository"
	"github.com/joshbatley/proxy/utils"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config, err := utils.LoadConfig("./config.yml")
	if err != nil {
		panic("Config unreadable")
	}

	// DB setup
	db := database.Conn()
	c := repository.CacheRepository{
		Database: db,
	}

	q := handler.QueryHandler{
		CacheRepository: &c,
	}

	r := mux.NewRouter()
	r.HandleFunc("/config", handler.ClientServe)
	r.HandleFunc("/query", q.Serve)
	http.Handle("/", r)

	log.Println("Listing on localhosts:" + config.Port)
	http.ListenAndServe("localhost:"+config.Port, nil)
}
