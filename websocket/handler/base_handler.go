package handler

import (
	"net/http"

	"github.com/uenoryo/chitoi/data"
	"golang.org/x/net/websocket"
)

// WriteJSON is XXX
func WriteJSON(ws *websocket.Conn, res interface{}) {
	websocket.JSON.Send(ws, res)
}

// WriteSuccess is XXX
func WriteSuccess(ws *websocket.Conn, res interface{}) {
	WriteJSON(ws, data.BaseResponse{
		Code:         http.StatusOK,
		DebugMessage: "OK",
		Message:      "",
		Data:         res,
	})
}

// WriteError is XXX
func WriteError(ws *websocket.Conn, code int, msg, debugMsg string) {
	WriteJSON(ws, data.BaseResponse{
		Code:         code,
		DebugMessage: debugMsg,
		Message:      msg,
		Data:         []string{},
	})
}

// WriteError500 is XXX
func WriteError500(ws *websocket.Conn, debugMsg string) {
	WriteError(ws, http.StatusInternalServerError, "エラーが発生しました", debugMsg)
}

// WriteError400 is XXX
func WriteError400(ws *websocket.Conn, debugMsg string) {
	WriteError(ws, http.StatusBadRequest, "無効なリクエストです", debugMsg)
}

// WriteError404 is XXX
func WriteError404(ws *websocket.Conn) {
	WriteError(ws, http.StatusNotFound, "ページが見つかりません", "ページが見つかりません")
}

// WriteError400or500 is XXX
func WriteError400or500(ws *websocket.Conn, err error) {
	switch e := err.(type) {
	case *data.ErrorBadRequest:
		WriteError400(ws, e.Error())
	default:
		WriteError500(ws, err.Error())
	}
}
