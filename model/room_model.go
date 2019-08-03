package model

import (
	"github.com/uenoryo/chitoi/core"
	"github.com/uenoryo/chitoi/database/row"
)

type Room struct {
	Row  *row.Room
	core *core.Core
}
