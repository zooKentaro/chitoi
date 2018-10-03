package api

import (
	"net/http"

	"github.com/uenoryo/chitoi/handler"
)

// NewServer is XXX
func NewServer() *http.ServeMux {
	server := http.NewServeMux()

	userHandler := handler.NewUserServer(handler.NewUserHandler())

	server.Handle("/user/", userHandler)

	return server
}
