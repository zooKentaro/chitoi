package handler

import (
	"log"
	"net/http"

	"github.com/uenoryo/chitoi/data"
	"github.com/uenoryo/chitoi/service"
)

// NewBusinessHandler is XXX
func NewBusinessHandler(srv service.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		Service: srv,
	}
}

type BusinessHandler struct {
	Service service.BusinessService
}

// NewBusinessServer is XXX
func NewBusinessServer(h *BusinessHandler) *http.ServeMux {
	server := http.NewServeMux()
	server.HandleFunc("/business/list", h.ListHandler)

	return server
}

// ListHandler is XXX
func (h *BusinessHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	req := &data.BusinessListRequest{}
	err := ScanRequest(r, req)
	if err != nil {
		log.Println(err.Error())
		WriteError400(w, err.Error())
		return
	}

	res, err := h.Service.List(req)
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
