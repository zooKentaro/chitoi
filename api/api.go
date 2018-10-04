package api

import (
	"net/http"

	"github.com/uenoryo/chitoi/handler"
	"github.com/uenoryo/chitoi/service"
)

// NewServer is XXX
func NewServer() *http.ServeMux {
	server := http.NewServeMux()

	userService := service.NewUserService()
	userHandler := handler.NewUserServer(handler.NewUserHandler(userService))

	server.Handle("/user/", userHandler)

	return server
}
