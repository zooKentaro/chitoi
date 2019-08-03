package model

import (
    "time"

    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/database/row"
    "github.com/uenoryo/chitoi/lib/id"
)

const (
    InsertRoomSQL = "INSERT INTO `room` (`id`, `owner_id`, `user1_id`, `user2_id`, `user3_id`, `user4_id`) VALUES (?,?,?,?,?)"
)

type UserRoom struct {
    user *User
    core *core.Core
}

func (ur *UserRoom) Create() (*Room, error) {
    room := &row.Room{
        ID:        id.Generate(),
        OwnerID:   ur.user.Row.ID,
        User1ID:   ur.user.Row.ID,
        User2ID:   0,
        User3ID:   0,
        User4ID:   0,
        CreatedAt: time.Now(),
    }
    if _, err := ur.core.DB.Exec(InsertRoomSQL, room.ID, room.OwnerID, room.User1ID, room.User2ID, room.User3ID, room.User4ID, room.CreatedAt); err != nil {
        return nil, errors.Wrapf(err, "error create room, sql:%s", InsertRoomSQL)
    }

    return &Room{
        Row:  room,
        core: ur.core,
    }, nil
}
