package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/database"
	"github.com/joshbatley/proxy/internal/handler"
	"github.com/joshbatley/proxy/internal/store"
	"github.com/joshbatley/proxy/internal/utils"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config, err := utils.LoadConfig("./config.yml")
	if err != nil {
		panic("Config unreadable")
	}
	// DB setup
	db := database.Conn()

	q := handler.QueryHandler{
		CacheRepository: &store.CacheRepository{
			Database: db,
		},
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
