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

func CreateNewRoom(core *core.Core, userID uint64) (*Room, error) {
    room := &row.Room{
        ID:        id.Generate(),
        OwnerID:   userID,
        User1ID:   userID,
        User2ID:   0,
        User3ID:   0,
        User4ID:   0,
        CreatedAt: time.Now(),
    }
    if _, err := core.DB.Exec(InsertRoomSQL, room.ID, room.OwnerID, room.User1ID, room.User2ID, room.User3ID, room.User4ID, room.CreatedAt); err != nil {
        return nil, errors.Wrapf(err, "error create room, sql:%s", InsertRoomSQL)
    }

    return &Room{
        Row:  room,
        core: core,
    }, nil
}

type Room struct {
    Row  *row.Room
    core *core.Core
}
