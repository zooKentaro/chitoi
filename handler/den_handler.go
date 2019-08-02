package handler

import (
	"net/http"

	"github.com/uenoryo/chitoi/core"
	"github.com/uenoryo/chitoi/websocket/service"
	"golang.org/x/net/websocket"
)

// NewDenHandler is XXX
func NewDenHandler(core *core.Core, srv service.DenService) *DenHandler {
	return &DenHandler{
		Core:    core,
		Service: srv,
	}
}

type DenHandler struct {
	Core    *core.Core
	Service service.DenService
}

// NewDenServer is XXX
func NewDenServer(h *DenHandler) *http.ServeMux {
	server := http.NewServeMux()
	server.Handle("/entry", websocket.Handler(h.EntryHandler))

	return server
}

// EntryHandler is XXX
func (h *DenHandler) EntryHandler(ws *websocket.Conn) {
	h.Service.Entry(ws)
}
