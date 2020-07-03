package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/client"
	"github.com/joshbatley/proxy/config"
)

func main() {
	config, err := config.Load("./config.yml")
	r := mux.NewRouter()

	if err != nil {
		panic("Config unreadable")
	}

	spa := client.SpaHandler{StaticPath: "./webapp/build", IndexPath: "index.html"}
	r.PathPrefix("/config").Handler(spa)

	http.Handle("/", r)
	log.Println("Listing on localhosts:" + config.Port)
	http.ListenAndServe("localhost:"+config.Port, nil)

}
