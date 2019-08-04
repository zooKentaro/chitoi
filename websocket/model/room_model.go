package model

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/uenoryo/chitoi/core"
	"github.com/uenoryo/chitoi/database/row"
	"github.com/uenoryo/chitoi/model"
)

const (
	FindRoomByCodeSQL = "SELECT * FROM room WHERE code = ? AND expired_at < ?"
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
