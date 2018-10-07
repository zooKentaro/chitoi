package main

import (
	"log"
	"net/http"

	"github.com/uenoryo/chitoi/api"
	"github.com/uenoryo/chitoi/env"
)

func main() {
	if err := env.Load(); err != nil {
		log.Fatal("error load env, error: ", err.Error())
		return
	}

	server, err := api.NewServer()
	if err != nil {
		log.Fatal("error new server, error: ", err.Error())
		return
	}

	srv := http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err.Error())
	}
}
