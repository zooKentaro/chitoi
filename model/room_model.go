package model

import (
	"github.com/uenoryo/chitoi/core"
	"github.com/uenoryo/chitoi/database/row"
)

// NewRoom (､´･ω･)▄︻┻┳═一
func NewRoom(core *core.Core, row *row.Room) *Room {
	return &Room{
		core,
		row,
	}
}

// Room はDenを行う部屋
type Room struct {
	core *core.Core
	Row  *row.Room
}
