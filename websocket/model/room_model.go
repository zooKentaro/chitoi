package model

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/uenoryo/chitoi/core"
	"github.com/uenoryo/chitoi/database/row"
	"github.com/uenoryo/chitoi/model"
	"github.com/uenoryo/chitoi/websocket/packet"
	"golang.org/x/net/websocket"
)

const (
	FindRoomByCodeSQL   = "SELECT * FROM room WHERE code = ? AND expired_at > ?"
	UpdateRoomByCodeSQL = "UPDATE room SET player1_id = ?, player2_id = ?, player3_id = ?, player4_id = ?, expired_at = ? WHERE code = ?"
)

var (
	ErrRoomHasNoVacancie = errors.New("error room has no vacancie")
)

// NewRoomRepository (､´･ω･)▄︻┻┳═一
func NewRoomRepository(core *core.Core) *RoomRepository {
	return &RoomRepository{core}
}

type RoomRepository struct {
	core *core.Core
}

func (repo *RoomRepository) FindByCode(code uint32) (*Room, error) {
	row := row.Room{}
	err := repo.core.DB.Get(&row, FindRoomByCodeSQL, code, time.Now())
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Errorf("room code:%d is not found", code)
	case err != nil:
		return nil, errors.Wrapf(err, "error find room by code, sql:%s", FindRoomByCodeSQL)
	}

	users, err := NewUserRepository(repo.core).FindByIDs(row.PlayerIDs())
	if err != nil {
		return nil, errors.Wrap(err, "find room player failed")
	}
	userByID := make(map[uint64]*User, len(users))
	for _, u := range users {
		userByID[u.Row.ID] = u
	}
	return NewRoom(repo.core, &row, userByID[row.Player1ID], userByID[row.Player2ID]), nil
}

// Save (､´･ω･)▄︻┻┳═一
func (repo *RoomRepository) Save(room *Room) error {
	if _, err := repo.core.DB.Exec(UpdateRoomByCodeSQL, room.Row.Player1ID, room.Row.Player2ID, room.Row.Player3ID, room.Row.Player4ID, time.Now().Add(time.Hour*24*2), room.Row.Code); err != nil {
		return errors.Wrapf(err, "error update room, sql:%s", UpdateRoomByCodeSQL)
	}
	return nil
}

type Room struct {
	core           *core.Core
	Row            *row.Room
	server         *Server
	Clients        map[uint64]*Client
	authentication map[string]uint64
	Player1        *User
	Player2        *User
}

func NewRoom(core *core.Core, row *row.Room, player1, player2 *User) *Room {
	return &Room{
		core,
		row,
		nil,
		make(map[uint64]*Client),
		make(map[string]uint64),
		player1,
		player2,
	}
}

func (r *Room) RegisterServer(server *Server) {
	r.server = server
}

func (r *Room) Server() (*Server, error) {
	if r.server == nil {
		return nil, errors.Errorf("room is not registered server")
	}
	return r.server, nil
}

// OwnerIs は user が room のオーナーかどうかを返す
func (r *Room) OwnerIs(user *model.User) bool {
	return r.Row.OwnerID == user.Row.ID
}

// Entry はroomにuserを入室させる
func (r *Room) Entry(ws *websocket.Conn, user *User, sessionID string) error {
	switch {
	case r.isAlreadyEntried(user):
		break
	case r.Row.Player1ID == 0:
		r.Row.Player1ID = user.Row.ID
		r.Player1 = user
	case r.Row.Player2ID == 0:
		r.Row.Player2ID = user.Row.ID
		r.Player2 = user
	default:
		return ErrRoomHasNoVacancie
	}
	r.Clients[user.Row.ID] = NewClient(ws, r, user)

	// TODO: 本来はEntryした時にTokenを生成して返し、以降の通信ではそれを利用する形にする
	r.authentication[sessionID] = user.Row.ID
	return nil
}

func (r *Room) PushOut(client *Client) {
	delete(r.Clients, client.ID())
}

func (r *Room) ListenAllClients() {
	for _, client := range r.Clients {
		if client.IsListening {
			continue
		}
		client.Listen()
	}
}

func (r *Room) Authenticate(sessionID string) (uint64, error) {
	userID, ok := r.authentication[sessionID]
	if !ok {
		return 0, errors.New("unauthorized user")
	}
	return userID, nil
}

func (r *Room) Submit(request *packet.RequestPacket) error {
	request.RoomCode = r.Row.Code
	server, err := r.Server()
	if err != nil {
		return errors.Wrap(err, "error get server")
	}
	if err := server.Validate(request); err != nil {
		return errors.Wrap(err, "failed validate packet")
	}
	server.Receive(request)
	return nil
}

// SubmitOnExitPlayer は player の退室を伝えるためにbroadcastする
// 通信切断時にも呼ばれる => SessionID等受け取れないので認証できない
func (r *Room) SubmitOnExitPlayer(request *packet.RequestPacket) error {
	request.RoomCode = r.Row.Code
	request.Method = packet.MethodExitPlayer
	server, err := r.Server()
	if err != nil {
		return errors.Wrap(err, "error get server")
	}
	server.Receive(request)
	return nil
}

// SubmitOnEnterPlayer は player の入室を伝えるためにbroadcastする
// 認証後に呼ばれるので認証しない
func (r *Room) SubmitOnEnterPlayer() error {
	request := &packet.RequestPacket{
		RoomCode: r.Row.Code,
		Method:   packet.MethodEntryPlayer,
	}
	server, err := r.Server()
	if err != nil {
		return errors.Wrap(err, "error get server")
	}
	server.Receive(request)
	return nil
}

func (r *Room) SendToMembers(pkt *packet.BroadcastPacket) {
	for _, member := range r.Clients {
		member.Receive(pkt)
	}
}

func (r *Room) isAlreadyEntried(user *User) bool {
	switch user.Row.ID {
	case r.Row.Player1ID, r.Row.Player2ID:
		return true
	default:
		return false
	}
}
