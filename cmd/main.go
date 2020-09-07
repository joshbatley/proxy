package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/database"
	"github.com/joshbatley/proxy/internal/engine"
	"github.com/joshbatley/proxy/internal/handler"
	"github.com/joshbatley/proxy/internal/store"
	"github.com/joshbatley/proxy/internal/utils"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// TODO: Set up logger

	// TODO: Set up as flags
	config, err := utils.LoadConfig("./config.yml")
	if err != nil {
		panic("Config unreadable")
	}
	// DB setup
	db := database.Conn()
	store := &store.Store{
		Database: db,
	}

	q := handler.QueryHandler{
		Rules: engine.NewEngine(store),
		Store: store,
	}

	r := mux.NewRouter()
	r.SkipClean(true)
	r.UseEncodedPath()

	r.PathPrefix("/{config:config.*}").Handler(handler.ClientHandler{
		StaticPath: "./webapp/build",
		IndexPath:  "index.html",
	})
	r.PathPrefix("/{collection:[0-9]*}/{query:.*}").Handler(q)
	r.PathPrefix("/{query:.*}").Handler(q)

	log.Println("Listing on localhosts:" + config.Port)
	http.ListenAndServe("localhost:"+config.Port, r)
}
