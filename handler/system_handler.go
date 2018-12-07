package handler

import (
	"log"
	"net/http"

	"github.com/uenoryo/chitoi/core"
	"github.com/uenoryo/chitoi/data"
	"github.com/uenoryo/chitoi/service"
)

// NewSystemHandler is XXX
func NewSystemHandler(core *core.Core, srv service.GameService) *SystemHandler {
	return &SystemHandler{
		Core:    core,
		Service: srv,
	}
}

type SystemHandler struct {
	Core    *core.Core
	Service service.GameService
}

// NewSystemServer is XXX
func NewSystemServer(h *SystemHandler) *http.ServeMux {
	server := http.NewServeMux()
	server.HandleFunc("/system/logger", h.LoggerHandler)

	return server
}

// LoggerHandler はリクエストされた内容をそのままログに出力します
func (h *SystemHandler) LoggerHandler(w http.ResponseWriter, r *http.Request) {
	req := &data.SystemLoggerRequest{}
	err := ScanRequest(r, req)
	if err != nil {
		log.Println(err.Error())
		if err := h.Core.Logger.PostError("api.system.logger", err.Error()); err != nil {
			log.Println(err.Error())
		}
		WriteError400(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		log.Println(err.Error())
		if err := h.Core.Logger.PostError("api.system.logger", err.Error()); err != nil {
			log.Println(err.Error())
		}
		WriteError400(w, err.Error())
		return
	}

	if err := h.Core.Logger.PostError(req.Tag, req.Body); err != nil {
		log.Println(err.Error())
	}

	if err = WriteSuccess(w, data.SystemLoggerResponse{}); err != nil {
		log.Println(err.Error())
		return
	}
	return
}
