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
	SessionIDKey = "SESSION-ID"
	RoomCodeKey  = "ROOM-CODE"
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
	cnt     int
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
	sessionIDCookie, err := ws.Request().Cookie(SessionIDKey)
	if err != nil {
		return errors.Wrapf(err, "error get session id from cookie, key:%s", SessionIDKey)
	}
	sessionID := sessionIDCookie.Value

	roomCodeCookie, err := ws.Request().Cookie(RoomCodeKey)
	if err != nil {
		return errors.Wrapf(err, "error get room code from cookie, key:%s", RoomCodeKey)
	}

	roomCodeInt, err := strconv.Atoi(roomCodeCookie.Value)
	if err != nil {
		return errors.Errorf("invalid room code:%s", roomCodeCookie.Value)
	}
	roomCode := uint32(roomCodeInt)

	user, err := apiservice.NewAuthService(srv.core).Authenticate(sessionID)
	if err != nil {
		return errors.Wrap(err, "error authenticate user")
	}

	roomRepo := model.NewRoomRepository(srv.core)

	userRoom, err := roomRepo.FindByCode(roomCode)
	if err != nil {
		return errors.Wrapf(err, "error find room by code:%d", roomCode)
	}

	var room *model.Room
	if userRoom.OwnerIs(user) {
		room = srv._server.Launch(userRoom)
	} else {
		r := srv._server.LaunchedRoom(roomCode)
		if r == nil {
			return errors.Errorf("room:%s is not launched on server", roomCode)
		}
		room = r
	}

	if err := room.Entry(ws, user, sessionID); err != nil {
		return errors.Wrapf(err, "error entry room, room code:%s", roomCode)
	}
	if err := roomRepo.Save(room); err != nil {
		return errors.Wrap(err, "error save room")
	}

	room.ListenAllClients()
	return nil
}
