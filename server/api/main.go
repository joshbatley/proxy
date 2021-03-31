package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/server/api/admin"
	"github.com/joshbatley/proxy/server/api/client"
	"github.com/joshbatley/proxy/server/api/query"
	"github.com/joshbatley/proxy/server/domain/collections"
	"github.com/joshbatley/proxy/server/domain/endpoints"
	"github.com/joshbatley/proxy/server/domain/responses"
	"github.com/joshbatley/proxy/server/domain/rules"
	"github.com/joshbatley/proxy/server/internal/database"
	"github.com/joshbatley/proxy/server/internal/logger"
	"github.com/joshbatley/proxy/server/internal/migration"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log := logger.Setup()

	port := flag.String("port", "5000", "desired port for internal server to run on")
	flag.Parse()

	// Start Migratio
	err := migration.StartUp()
	if err != nil {
		log.Fatal(err)
	}

	// DB setup
	db, err := database.Conn()
	if err != nil {
		log.Fatal(err)
	}

	collections := collections.NewManager(collections.NewSQLRepository(db))
	endpoints := endpoints.NewManager(endpoints.NewSQLRepository(db))
	responses := responses.NewManager(responses.NewSQLRepository(db))
	rules := rules.NewManager(rules.NewSQLRepository(db))

	q := query.NewHandler(
		collections,
		endpoints,
		responses,
		rules,
		log,
	)
	a := admin.NewHandler(
		collections,
		endpoints,
		responses,
		rules,
		log,
	)

	r := mux.NewRouter().SkipClean(true).UseEncodedPath()

	adminRouter := r.PathPrefix("/admin").Subrouter()
	a.Router(adminRouter)

	r.PathPrefix("/config").Handler(client.Handler{
		StaticPath: "../build",
		IndexPath:  "index.html",
	})
	r.PathPrefix("/ping").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("pong")
		w.Write([]byte("pong"))
	})

	r.PathPrefix("/{collection:[0-9]*}/{query:.*}").Handler(q)
	r.PathPrefix("/{query:.*}").Handler(q)

	log.Infof("Listing on localhosts:" + *port)
	err = http.ListenAndServe(":"+*port, r)
	log.Info(err)
}
