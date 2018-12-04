package handler

import (
	"log"
	"net/http"

	"github.com/uenoryo/chitoi/core"
	"github.com/uenoryo/chitoi/data"
	"github.com/uenoryo/chitoi/service"
)

// NewBusinessHandler is XXX
func NewBusinessHandler(core *core.Core, srv service.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		Core:    core,
		Service: srv,
	}
}

type BusinessHandler struct {
	Core    *core.Core
	Service service.BusinessService
}

// NewBusinessServer is XXX
func NewBusinessServer(h *BusinessHandler) *http.ServeMux {
	server := http.NewServeMux()
	server.HandleFunc("/business/list", h.ListHandler)
	server.HandleFunc("/business/buy", h.BuyHandler)

	return server
}

// ListHandler is XXX
func (h *BusinessHandler) ListHandler(w http.ResponseWriter, r *http.Request) {
	res, err := h.Service.List()
	if err != nil {
		log.Println(err.Error())
		if err := h.Core.Logger.PostError("api.business.list", err.Error()); err != nil {
			log.Println(err.Error())
		}
		WriteError400or500(w, err)
		return
	}

	if err = WriteSuccess(w, res); err != nil {
		log.Println(err.Error())
		return
	}
	return
}

// BuyHandler is XXX
func (h *BusinessHandler) BuyHandler(w http.ResponseWriter, r *http.Request) {
	req := &data.BusinessBuyRequest{}
	err := ScanRequest(r, req)
	if err != nil {
		log.Println(err.Error())
		if err := h.Core.Logger.PostError("api.business.buy", err.Error()); err != nil {
			log.Println(err.Error())
		}
		WriteError400(w, err.Error())
		return
	}

	res, err := h.Service.Buy(req)
	if err != nil {
		log.Println(err.Error())
		if err := h.Core.Logger.PostError("api.business.buy", err.Error()); err != nil {
			log.Println(err.Error())
		}
		WriteError400or500(w, err)
		return
	}

	if err = WriteSuccess(w, res); err != nil {
		log.Println(err.Error())
		return
	}
	return
}
