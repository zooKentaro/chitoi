package handler

import (
	"fmt"
	"net/http"
)

// NewUserHandler is XXX
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

type UserHandler struct{}

// NewUserServer is XXX
func NewUserServer(h *UserHandler) *http.ServeMux {
	server := http.NewServeMux()
	server.HandleFunc("/user/signup", h.SignupHandler)

	return server
}

// SignupHandler is XXX
func (h *UserHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "user")
}
