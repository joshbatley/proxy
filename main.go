package main

import (
	"fmt"
	"net/http"

	"goproxy/config"
	"goproxy/query"
)

func main() {
	config, err := config.Load("./config.yml")

	if err != nil {
		panic("Config unreadable")
	}

	// Set up DB

	http.HandleFunc("/query", query.Serve)

	fmt.Println("listing on 127.0.0.1:" + config.Port)

	http.ListenAndServe(":"+config.Port, nil)
}
