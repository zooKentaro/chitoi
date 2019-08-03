package handler

import (
	"log"
	"net/http"

	"github.com/uenoryo/chitoi/core"
	"github.com/uenoryo/chitoi/data"
	"github.com/uenoryo/chitoi/service"
)

// NewRoomHandler is XXX
func NewRoomHandler(core *core.Core, srv service.RoomService) *RoomHandler {
	return &RoomHandler{
		Core:    core,
		Service: srv,
	}
}

type RoomHandler struct {
	Core    *core.Core
	Service service.RoomService
}

// NewRoomServer is XXX
func NewRoomServer(h *RoomHandler) *http.ServeMux {
	server := http.NewServeMux()
	server.HandleFunc("/room/create", h.CreateHandler)

	return server
}

// CreateHandler is XXX
func (h *RoomHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	req := &data.RoomCreateRequest{}
	if err := ScanRequest(r, req); err != nil {
		log.Println(err.Error())
		if err := h.Core.Logger.PostError("api.business.buy", err.Error()); err != nil {
			log.Println(err.Error())
		}
		WriteError400(w, err.Error())
		return
	}

	res, err := h.Service.Create(req)
	if err != nil {
		log.Println(err.Error())
		if err := h.Core.Logger.PostError("api.room.list", err.Error()); err != nil {
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
