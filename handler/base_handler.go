package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/uenoryo/chitoi/data"
)

// ScanRequest is XXX
func ScanRequest(r *http.Request, req interface{}) error {
	if r.Method != "POST" {
		return errors.Errorf("%s request is not allowed", r.Method)
	}

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
func WriteJSON(w http.ResponseWriter, res interface{}) error {
	response, err := json.Marshal(res)
	if err != nil {
		return err
	}

	_, err = w.Write(response)
	if err != nil {
		return err
	}

	return nil
}

func WriteBaseHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	if os.Getenv("CHITOI_ENV") != "production" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
}

// WriteSuccess is XXX
func WriteSuccess(w http.ResponseWriter, res interface{}) error {
	WriteBaseHeader(w)
	w.WriteHeader(http.StatusOK)

	return WriteJSON(w, data.BaseResponse{
		Code:         http.StatusOK,
		DebugMessage: "OK",
		Message:      "",
		Data:         res,
	})
}

// WriteError is XXX
func WriteError(w http.ResponseWriter, code int, msg, debugMsg string) error {
	WriteBaseHeader(w)
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

// WriteError400or500 is XXX
func WriteError400or500(w http.ResponseWriter, err error) error {
	switch e := err.(type) {
	case *data.ErrorBadRequest:
		return WriteError400(w, e.Error())
	}
	return WriteError500(w, err.Error())
}
