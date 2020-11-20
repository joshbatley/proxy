package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/api/handler/admin"
	"github.com/joshbatley/proxy/api/handler/client"
	"github.com/joshbatley/proxy/api/handler/query"
	"github.com/joshbatley/proxy/domain/collections"
	"github.com/joshbatley/proxy/domain/endpoints"
	"github.com/joshbatley/proxy/domain/responses"
	"github.com/joshbatley/proxy/domain/rules"
	"github.com/joshbatley/proxy/internal/database"
	"github.com/joshbatley/proxy/internal/logger"
	"github.com/joshbatley/proxy/internal/migration"
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

	r := mux.NewRouter().SkipClean(true).UseEncodedPath()

	adminRouter := r.PathPrefix("/admin").Subrouter()
	admin.Router(adminRouter)

	r.PathPrefix("/{config:config.*}").Handler(client.Handler{
		StaticPath: "./webapp/build",
		IndexPath:  "index.html",
	})
	r.PathPrefix("/{collection:[0-9]*}/{query:.*}").Handler(q)
	r.PathPrefix("/{query:.*}").Handler(q)

	log.Infof("Listing on localhosts:" + *port)
	err = http.ListenAndServe("localhost:"+*port, r)
	log.Info(err)
}
