package server

import (
    "net/http"

    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/handler"
    "github.com/uenoryo/chitoi/service"
)

// NewServer is XXX
func NewServer() (*http.ServeMux, *service.Listener, error) {
    srv := http.NewServeMux()

    core, err := core.New()
    if err != nil {
        return nil, nil, errors.Wrap(err, "error new core")
    }
    if err := core.LoadMasterdata(); err != nil {
        return nil, nil, errors.Wrap(err, "error load masterdata")
    }

    denService := service.NewDenService(core)
    listener := denService.Listener()
    denHandler := handler.NewDenServer(handler.NewDenHandler(core, denService))

    srv.Handle("/den/", denHandler)

    return srv, listener, nil
}
