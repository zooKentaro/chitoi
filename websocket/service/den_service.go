package service

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/uenoryo/chitoi/core"
	apiservice "github.com/uenoryo/chitoi/service"
	"github.com/uenoryo/chitoi/websocket/model"
	"golang.org/x/net/websocket"
)

const (
	UserSessionHeaderKey = "X-CHITOI-SESSION"
	RoomCodeHeaderKey    = "X-CHITOI-ROOM-CODE"
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
	server := model.NewServer(core)

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
	var (
		sessionID   = ws.Request().Header.Get(UserSessionHeaderKey)
		roomCodeStr = ws.Request().Header.Get(RoomCodeHeaderKey)
	)
	roomCodeInt, err := strconv.Atoi(roomCodeStr)
	if err != nil {
		return errors.Wrapf(err, "invalid room code:%s", roomCodeStr)
	}
	roomCode := uint32(roomCodeInt)

	user, err := apiservice.NewAuthService(srv.core).Authenticate(sessionID)
	if err != nil {
		return errors.Wrap(err, "error authenticate user")
	}

	room, err := model.NewRoomRepository(srv.core).FindByCode(roomCode)
	if err != nil {
		return errors.Wrapf(err, "error find room by code:%s", roomCode)
	}

	if room.OwnerIs(user) {
		srv._server.Launch(room)
	}

	if !srv._server.IsLaunched(room) {
		return errors.Errorf("room:%s is not launched on server", roomCode)
	}

	if err := room.Entry(user); err != nil {
		return errors.Wrapf(err, "error entry room, room code:%s", roomCode)
	}
	return nil
}
