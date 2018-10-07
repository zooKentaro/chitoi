package handler

import (
	"log"
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
	service := service.NewUserService()
	res, err := service.Signup()
	if err != nil {
		log.Println(err.Error())
		WriteError500(w, err.Error())
	}

	if err = WriteJSON(w, res); err != nil {
		log.Println(err.Error())
	}
	return
}
