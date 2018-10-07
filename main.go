package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/uenoryo/chitoi/api"
)

func main() {
	if os.Getenv("CHITOI_ENV") == "" {
		os.Setenv("CHITOI_ENV", "development")
	}

	err := godotenv.Load(fmt.Sprintf(".env.%s", os.Getenv("CHITOI_ENV")))
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	server := api.NewServer()

	srv := http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err.Error())
	}
}
