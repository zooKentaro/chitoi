package handler

import (
	"net/http"

	"github.com/uenoryo/chitoi/core"
	"github.com/uenoryo/chitoi/service"
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
	server.Handle("/", websocket.Handler(h.Handler))

	return server
}

// DenHandler is XXX
func (h *DenHandler) Handler(ws *websocket.Conn) {
	h.Service.Entry(ws)
}
