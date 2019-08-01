package server

import (
    "log"
    "net/http"

    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/handler"
    "github.com/uenoryo/chitoi/service"
)

// NewServer is XXX
func NewServer() (*http.ServeMux, *Listener, error) {
    srv := http.NewServeMux()

    core, err := core.New()
    if err != nil {
        return nil, nil, errors.Wrap(err, "error new core")
    }
    if err := core.LoadMasterdata(); err != nil {
        return nil, nil, errors.Wrap(err, "error load masterdata")
    }

    denService := service.NewUserService(core)
    denHandler := handler.NewDenServer(handler.NewDenHandler(core, denService))

    srv.Handle("/den/", denHandler)

    listener := newListener()

    return srv, listener, nil
}

func newListener() *Listener {
    var (
        doneCh = make(chan bool)
        errCh  = make(chan error)
    )
    return &Listener{
        doneCh,
        errCh,
    }
}

type Listener struct {
    doneCh chan bool
    errCh  chan error
}

func (l *Listener) Listen() {
    for {
        select {
        case err := <-l.errCh:
            log.Println("Error:", err.Error())

        case <-l.doneCh:
            return
        }
    }
}
