package model

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/uenoryo/chitoi/core"
	"github.com/uenoryo/chitoi/database/row"
	"github.com/uenoryo/chitoi/model"
	"golang.org/x/net/websocket"
)

const (
	FindRoomByCodeSQL   = "SELECT * FROM room WHERE code = ? AND expired_at < ?"
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
	default:
		return NewRoom(repo.core, &row), nil
	}
}

// Save (､´･ω･)▄︻┻┳═一
func (repo *RoomRepository) Save(room *Room) error {
	if _, err := repo.core.DB.Exec(UpdateRoomByCodeSQL, room.Row.Player1ID, room.Row.Player2ID, room.Row.Player3ID, room.Row.Player4ID, time.Now(), room.Row.Code); err != nil {
		return errors.Wrapf(err, "error update room, sql:%s", UpdateRoomByCodeSQL)
	}
	return nil
}

type Room struct {
	core    *core.Core
	Row     *row.Room
	Clients map[uint64]*Client
}

func NewRoom(core *core.Core, row *row.Room) *Room {
	return &Room{
		core,
		row,
		make(map[uint64]*Client),
	}
}

// OwnerIs は user が room のオーナーかどうかを返す
func (r *Room) OwnerIs(user *model.User) bool {
	return r.Row.OwnerID == user.Row.ID
}

// Entry はroomにuserを入室させる
func (r *Room) Entry(ws *websocket.Conn, user *model.User) error {
	switch {
	case r.isAlreadyEntried(user):
		return nil
	case r.Row.Player1ID == 0:
		r.Row.Player1ID = user.Row.ID
	case r.Row.Player2ID == 0:
		r.Row.Player2ID = user.Row.ID
	default:
		return ErrRoomHasNoVacancie
	}
	r.Clients[user.Row.ID] = NewClient(ws, r, user)

	return nil
}

func (r *Room) ListenAllClients() {
	for _, client := range r.Clients {
		if client.IsListening {
			continue
		}
		client.Listen()
	}
}

func (r *Room) isAlreadyEntried(user *model.User) bool {
	switch user.Row.ID {
	case r.Row.Player1ID, r.Row.Player2ID:
		return true
	default:
		return false
	}
}
