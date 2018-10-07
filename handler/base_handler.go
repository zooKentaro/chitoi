package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/uenoryo/chitoi/data"
)

// ScanRequest is XXX
func ScanRequest(r *http.Request, req interface{}) error {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, req)
	if err != nil {
		return err
	}

	return nil
}

// WriteJSON is XXX
func WriteJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write(response)
	if err != nil {
		return err
	}

	return nil
}

// WriteError is XXX
func WriteError(w http.ResponseWriter, code int, msg, debugMsg string) error {
	w.WriteHeader(code)

	return WriteJSON(w, data.BaseResponse{
		Code:         code,
		DebugMessage: debugMsg,
		Message:      msg,
		Data:         []string{},
	})
}

// WriteError500 is XXX
func WriteError500(w http.ResponseWriter, debugMsg string) error {
	return WriteError(w, http.StatusInternalServerError, "エラーが発生しました", debugMsg)
}

// WriteError400 is XXX
func WriteError400(w http.ResponseWriter, debugMsg string) error {
	return WriteError(w, http.StatusBadRequest, "無効なリクエストです", debugMsg)
}

// WriteError404 is XXX
func WriteError404(w http.ResponseWriter) error {
	return WriteError(w, http.StatusNotFound, "ページが見つかりません", "ページが見つかりません")
}
