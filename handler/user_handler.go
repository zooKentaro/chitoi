package handler

import (
	"net/http"

	"github.com/uenoryo/chitoi/service"
)

// NewUserHandler is XXX
func NewUserHandler(srv service.UserService) *UserHandler {
	return &UserHandler{
		Service: srv,
	}
}

type UserHandler struct {
	Service service.UserService
}

// NewUserServer is XXX
func NewUserServer(h *UserHandler) *http.ServeMux {
	server := http.NewServeMux()
	server.HandleFunc("/user/signup", h.SignupHandler)

	return server
}

// SignupHandler is XXX
func (h *UserHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	WriteError404(w)
}
