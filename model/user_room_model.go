package model

import (
    "math/rand"
    "time"

    "github.com/Songmu/retry"
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/database/row"
    "github.com/uenoryo/chitoi/lib/id"
)

const (
    InsertRoomSQL           = "INSERT INTO `room` (`id`, `owner_id`, `user1_id`, `user2_id`, `user3_id`, `user4_id`, `created_at`) VALUES (?,?,?,?,?,?,?)"
    CountValidRoomByCodeSQL = "SELECT count(*) FROM room WHERE code = ? AND expired_at > ?"
)

type UserRoom struct {
    core *core.Core
    user *User
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

func (ur *UserRoom) generateRoomCode() (uint32, error) {
    var (
        retryCount    = uint(10)
        retryInterval = time.Second * 0
        maxRandNum    = 900000
        minRandNum    = 100000
        result        = uint32(0)
    )
    err := retry.Retry(retryCount, retryInterval, func() error {
        rand.Seed(time.Now().UnixNano())
        code := uint32(rand.Intn(maxRandNum) + minRandNum)

        var count int
        if err := ur.core.DB.Select(&count, CountValidRoomByCodeSQL, code, time.Now()); err != nil {
            return errors.Wrapf(err, "count room failed, sql:%s", CountValidRoomByCodeSQL)
        }
        if count != 0 {
            return errors.Errorf("error room code:%d is already exists", code)
        }
        result = code
        return nil
    })
    if err != nil {
        return 0, errors.Wrap(err, "error generate room code")
    }
    return result, nil
}
