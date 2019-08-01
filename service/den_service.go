package service

import (
	"log"

	"github.com/uenoryo/chitoi/core"
)

type DenService interface {
	Listener() *Listener
}

type denService struct {
	core   *core.Core
	server *Listener
}

// NewDenService (､´･ω･)▄︻┻┳═一
func NewDenService(core *core.Core) DenService {
	listener := newListener()

	return &denService{
		core:   core,
		server: listener,
	}
}

func (srv *denService) Listener() *Listener {
	return srv.server
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
