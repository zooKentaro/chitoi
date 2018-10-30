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
	server.HandleFunc("/user/login", h.LoginHandler)
	server.HandleFunc("/user/info", h.InfoHandler)
	server.HandleFunc("/user/record", h.RecordHandler)

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

	if err := req.Validate(); err != nil {
		WriteError400(w, err.Error())
		return
	}

	res, err := h.Service.Signup(req)
	if err != nil {
		log.Println(err.Error())
		WriteError400or500(w, err)
		return
	}

	if err = WriteSuccess(w, res); err != nil {
		log.Println(err.Error())
		return
	}
	return
}

// LoginHandler is XXX
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	req := &data.UserLoginRequest{}
	err := ScanRequest(r, req)
	if err != nil {
		log.Println(err.Error())
		WriteError400(w, err.Error())
		return
	}

	res, err := h.Service.Login(req)
	if err != nil {
		log.Println(err.Error())
		WriteError400or500(w, err)
		return
	}

	if err = WriteSuccess(w, res); err != nil {
		log.Println(err.Error())
		return
	}
	return
}

// InfoHandler is XXX
func (h *UserHandler) InfoHandler(w http.ResponseWriter, r *http.Request) {
	req := &data.UserInfoRequest{}
	err := ScanRequest(r, req)
	if err != nil {
		log.Println(err.Error())
		WriteError400(w, err.Error())
		return
	}

	res, err := h.Service.Info(req)
	if err != nil {
		log.Println(err.Error())
		WriteError400or500(w, err)
		return
	}

	if err = WriteSuccess(w, res); err != nil {
		log.Println(err.Error())
		return
	}
	return
}

// RecordHandler is XXX
func (h *UserHandler) RecordHandler(w http.ResponseWriter, r *http.Request) {
	req := &data.UserRecordRequest{}
	err := ScanRequest(r, req)
	if err != nil {
		log.Println(err.Error())
		WriteError400(w, err.Error())
		return
	}

	res, err := h.Service.Record(req)
	if err != nil {
		log.Println(err.Error())
		WriteError400or500(w, err)
		return
	}

	if err = WriteSuccess(w, res); err != nil {
		log.Println(err.Error())
		return
	}
	return
}
