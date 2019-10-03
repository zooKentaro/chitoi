package model

import (
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
)

const (
    lockForUpdateUserSQL      = "SELECT * FROM user WHERE id = ? FOR UPDATE"
    updateUserByFinishGameSQL = "UPDATE user SET stamina = ?, money = ? WHERE id = ?"
)

// UserGameBehavior (､´･ω･)▄︻┻┳═一
type UserGameBehavior struct {
    core *core.Core
    user *User
}

// GameData は1ゲームのデータを扱う
type GameData struct {
    Money int64
}

// GameFinish は1ゲーム終了時の動作を行う
func (bvr *UserGameBehavior) Finish(data *GameData) error {
    if err := bvr.user.exhaustStamina(); err != nil {
        return errors.Wrap(err, "error exhaust stamina")
    }

    bvr.user.getOrLoseMoney(data.Money)

    if _, err := bvr.core.DB.Exec(lockForUpdateUserSQL, bvr.user.Row.ID); err != nil {
        return errors.Wrapf(err, "error lock for update, sql:%s", LockForUpdateUserSQL)
    }

    if _, err := bvr.core.DB.Exec(UpdateUserByFinishGameSQL, bvr.user.Row.Stamina, bvr.user.Row.Money, bvr.user.Row.ID); err != nil {
        return errors.Wrapf(err, "error update user data, sql:%s", updateUserByFinishGameSQL)
    }

    return nil
}
