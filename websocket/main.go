package main

import (
	"log"
	"net/http"

	"github.com/uenoryo/hnk/env"
)

func main() {
	if err := env.Load(); err != nil {
		log.Fatal("error load env, error: ", err.Error())
		return
	}

	server, listener, err := NewServer()
	if err != nil {
		log.Fatal("error new server, error: ", err.Error())
		return
	}

	go listener.Listen()

	srv := http.Server{
		Addr:    ":8081",
		Handler: server,
	}

	log.Println("web socket server is started.")
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err.Error())
	}
}
