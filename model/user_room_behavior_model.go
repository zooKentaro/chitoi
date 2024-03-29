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
    InsertRoomSQL           = "INSERT INTO `room` (`id`, `code`, `owner_id`, `player1_id`, `player2_id`, `player3_id`, `player4_id`, `created_at`, `expired_at`) VALUES (?,?,?,?,?,?,?,?,?)"
    CleanRoomSQL            = "DELETE FROM `room` WHERE owner_id = ?"
    CountValidRoomByCodeSQL = "SELECT count(*) FROM room WHERE code = ? AND expired_at > ?"
)

type UserRoomBehavior struct {
    core *core.Core
    user *User
}

func (ur *UserRoomBehavior) Create() (*Room, error) {
    roomCode, err := ur.generateRoomCode()
    if err != nil {
        return nil, errors.Wrap(err, "generate room code failed")
    }

    var (
        now       = time.Now()
        expiredAt = now.Add(time.Hour * 24 * 2)
    )
    room := &row.Room{
        ID:        id.Generate(),
        Code:      roomCode,
        OwnerID:   ur.user.Row.ID,
        Player1ID: ur.user.Row.ID,
        Player2ID: 0,
        Player3ID: 0,
        Player4ID: 0,
        CreatedAt: now,
        ExpiredAt: expiredAt,
    }
    if _, err := ur.core.DB.Exec(InsertRoomSQL, room.ID, room.Code, room.OwnerID, room.Player1ID, room.Player2ID, room.Player3ID, room.Player4ID, room.CreatedAt, room.ExpiredAt); err != nil {
        return nil, errors.Wrapf(err, "error create room, sql:%s", InsertRoomSQL)
    }

    return &Room{
        core: ur.core,
        Row:  room,
    }, nil
}

// Clean は自身がオーナーである部屋を一掃する
func (ur *UserRoomBehavior) Clean() error {
    if _, err := ur.core.DB.Exec(CleanRoomSQL, ur.user.Row.ID); err != nil {
        return errors.Wrapf(err, "error delete room, sql:%s", CleanRoomSQL)
    }
    return nil
}

func (ur *UserRoomBehavior) generateRoomCode() (uint32, error) {
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
        if err := ur.core.DB.Get(&count, CountValidRoomByCodeSQL, code, time.Now()); err != nil {
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
