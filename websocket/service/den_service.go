package service

import (
	"github.com/uenoryo/chitoi/core"
	"github.com/uenoryo/chitoi/websocket/model"
	"golang.org/x/net/websocket"
)

type DenService interface {
	Listener() Listener
	Entry(*websocket.Conn)
}

type Listener interface {
	Listen()
}

type denService struct {
	core    *core.Core
	_server *model.Server
}

// NewDenService (､´･ω･)▄︻┻┳═一
func NewDenService(core *core.Core) DenService {
	server := model.NewServer()

	return &denService{
		core:    core,
		_server: server,
	}
}

// Listener は listener を返す
func (srv *denService) Listener() Listener {
	return srv._server
}

// Entry は client を作成し、websocketで server と接続する
func (srv *denService) Entry(ws *websocket.Conn) {
	client := model.NewClient(ws, srv._server)
	srv._server.Add(client)
	client.Listen()
}
