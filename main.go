package main

import (
	"log"
	"net/http"

	"github.com/uenoryo/chitoi/api"
)

func main() {

	server := api.NewServer()

	srv := http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err.Error())
	}
}
