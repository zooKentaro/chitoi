package handler

import (
	"log"
	"net/http"

	"github.com/uenoryo/chitoi/data"
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
	req := &data.UserSignupRequest{}
	err := ScanRequest(r, req)
	if err != nil {
		log.Println(err.Error())
		WriteError400(w, err.Error())
		return
	}

	service := service.NewUserService()
	res, err := service.Signup(req)
	if err != nil {
		log.Println(err.Error())
		WriteError400or500(w, err)
	}

	if err = WriteJSON(w, res); err != nil {
		log.Println(err.Error())
	}
	return
}
