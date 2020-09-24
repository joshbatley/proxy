package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/api/handler"
	"github.com/joshbatley/proxy/domain/collections"
	"github.com/joshbatley/proxy/domain/endpoints"
	"github.com/joshbatley/proxy/domain/responses"
	"github.com/joshbatley/proxy/domain/rules"
	"github.com/joshbatley/proxy/internal/config"
	"github.com/joshbatley/proxy/internal/database"
	"github.com/joshbatley/proxy/internal/engine"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// TODO: Set up logger

	// TODO: Set up as flags
	config, err := config.Load("./config.yml")
	if err != nil {
		panic("Config unreadable")
	}
	// DB setup
	db := database.Conn()

	col := collections.NewManager(collections.NewSQLRepository(db))
	end := endpoints.NewManager(endpoints.NewSQLRepository(db))
	res := responses.NewManager(responses.NewSQLRepository(db))
	rul := rules.NewManager(rules.NewSQLRepository(db))
	eng := engine.NewEngine(rul, col)

	q := handler.NewQueryHandler(
		col,
		end,
		res,
		rul,
		eng,
	)

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
