package handler

import (
	"log"
	"net/http"

	"github.com/uenoryo/chitoi/core"
	"github.com/uenoryo/chitoi/data"
	"github.com/uenoryo/chitoi/service"
)

// NewGameHandler is XXX
func NewGameHandler(core *core.Core, srv service.GameService) *GameHandler {
	return &GameHandler{
		Core:    core,
		Service: srv,
	}
}

type GameHandler struct {
	Core    *core.Core
	Service service.GameService
}

// NewGameServer is XXX
func NewGameServer(h *GameHandler) *http.ServeMux {
	server := http.NewServeMux()
	server.HandleFunc("/game/finish", h.FinishHandler)

	return server
}

// FinishHandler is XXX
func (h *GameHandler) FinishHandler(w http.ResponseWriter, r *http.Request) {
	req := &data.GameFinishRequest{}
	err := ScanRequest(r, req)
	if err != nil {
		log.Println(err.Error())
		if err := h.Core.Logger.PostError("api.business", err.Error()); err != nil {
			log.Println(err.Error())
		}
		WriteError400(w, err.Error())
		return
	}

	res, err := h.Service.Finish(req)
	if err != nil {
		log.Println(err.Error())
		if err := h.Core.Logger.PostError("api.business", err.Error()); err != nil {
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
