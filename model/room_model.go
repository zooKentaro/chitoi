package model

import (
    "time"

    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/database/row"
)

const (
    InsertRoomSQL = "INSERT INTO `room` (`owner_id`, `user1_id`, `user2_id`, `user3_id`, `user4_id`) VALUES (?,?,?,?,?)"
)

func CreateNewRoom(core *core.Core, userID uint64) (*Room, error) {
    now := time.Now()

    roomRow := &row.Room{
        OwnerID:   userID,
        User1ID:   userID,
        User2ID:   0,
        User3ID:   0,
        User4ID:   0,
        CreatedAt: now,
    }
    res, err := core.DB.Exec(InsertRoomSQL, roomRow.OwnerID, roomRow.User1ID, roomRow.User2ID, roomRow.User3ID, roomRow.User4ID, roomRow.CreatedAt)
    if err != nil {
        return nil, errors.Wrapf(err, "error create room, sql:%s", InsertRoomSQL)
    }

    id, err := res.LastInsertId()
    if err != nil {
        return nil, errors.Wrap(err, "error last insert id")
    }
    roomRow.ID = uint64(id)

    return &Room{
        Row:  roomRow,
        core: core,
    }, nil
}

type Room struct {
    Row  *row.Room
    core *core.Core
}
