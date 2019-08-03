package service

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/uenoryo/chitoi/core"
	apiservice "github.com/uenoryo/chitoi/service"
	"github.com/uenoryo/chitoi/websocket/model"
	"golang.org/x/net/websocket"
)

const (
	UserSessionHeaderKey = "X-CHITOI-SESSION"
)

type DenService interface {
	Listener() Listener
	Entry(*websocket.Conn) error
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
func (srv *denService) Entry(ws *websocket.Conn) error {
	sessionID := ws.Request().Header.Get(UserSessionHeaderKey)

	user, err := apiservice.NewAuthService(srv.core).Authenticate(sessionID)
	if err != nil {
		return errors.Wrap(err, "error authenticate user")
	}
	fmt.Println(user)

	client := model.NewClient(ws, srv._server)
	srv._server.Add(client)
	client.Listen()
	return nil
}
