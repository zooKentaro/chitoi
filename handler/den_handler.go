package handler

import (
	"io"
	"net/http"

	"github.com/uenoryo/chitoi/core"
	"golang.org/x/net/websocket"
)

// NewDenHandler is XXX
func NewDenHandler(core *core.Core) *DenHandler {
	return &DenHandler{
		Core: core,
	}
}

type DenHandler struct {
	Core *core.Core
}

// NewDenServer is XXX
func NewDenServer(h *DenHandler) *http.ServeMux {
	server := http.NewServeMux()
	server.Handle("/", websocket.Handler(h.Handler))

	return server
}

// DenHandler is XXX
func (h *DenHandler) Handler(ws *websocket.Conn) {
	io.Copy(ws, ws)
}
