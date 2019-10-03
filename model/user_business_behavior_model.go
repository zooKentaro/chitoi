package model

import (
    "database/sql"
    "time"

    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/constant"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/database/row"
    "github.com/uenoryo/chitoi/lib/helper"
)

type UserBusinessBehavior struct {
    core *core.Core
    user *User
}

func (bvr *UserBusinessBehavior) BusinessBuy(business *Business) error {
    if err := business.IsOpen(); err != nil {
        return errors.Wrap(err, "error business is open")
    }

    if _, err := bvr.core.DB.Exec("SELECT * FROM user WHERE id = ? FOR UPDATE", bvr.user.Row.ID); err != nil {
        return errors.Wrap(err, "error lock for update")
    }

    exists := true
    ubRow := &row.UserBusiness{}
    q := "SELECT * FROM user_business WHERE user_id = ? AND business_id = ?"
    if err := bvr.core.DB.Get(ubRow, q, bvr.user.Row.ID, business.Row.ID); err != nil {
        if err == sql.ErrNoRows {
            exists = false
        } else {
            return errors.Wrap(err, "error find user business")
        }
    }

    if err := bvr.canBuy(ubRow, business); err != nil {
        return errors.Wrap(err, "error can buy")
    }

    nextPrice, err := business.NextPrice(ubRow)
    if err != nil {
        return errors.Wrap(err, "error next price")
    }
    bvr.user.spendMoney(nextPrice)

    if _, err := bvr.core.DB.Exec("UPDATE user SET money = ? WHERE id = ?", bvr.user.Row.Money, bvr.user.Row.ID); err != nil {
        return errors.Wrap(err, "error update user data")
    }

    if exists {
        q := "UPDATE user_business SET level = ?, last_buy_at = ? WHERE user_id = ? AND business_id = ?"
        if _, err := bvr.core.DB.Exec(q, ubRow.Level+1, time.Now(), bvr.user.Row.ID, business.Row.ID); err != nil {
            return errors.Wrap(err, "error update user data")
        }
    } else {
        q := "INSERT INTO user_business (user_id, business_id, level, last_buy_at) VALUES (?,?,?,?)"
        if _, err := bvr.core.DB.Exec(q, bvr.user.Row.ID, business.Row.ID, 1, time.Now()); err != nil {
            return errors.Wrap(err, "error create user business")
        }
    }
    return nil
}

// canBuy は購入可能な情状態かどうか以下のチェックを行います
// * 今日はまだ1度も購入していないこと
// * まだ Business Level が最大になっていないこと
// * 購入に必要な資金を持っていること
func (bvr *UserBusinessBehavior) canBuy(ub *row.UserBusiness, business *Business) error {
    if ub != nil {
        if helper.IsSameDay(ub.LastBuyAt, time.Now()) {
            return errors.New("cannot buy 2 times per day")
        }
        if ub.Level >= constant.MaxBusinessLevel {
            return errors.New("cannot level up more than this")
        }
    }

    nextPrice, err := business.NextPrice(ub)
    if err != nil {
        return errors.Wrap(err, "error next price")
    }

    if bvr.user.Row.Money < int64(nextPrice) {
        return errors.Errorf("error money is not enough, want %d but current %d", nextPrice, bvr.user.Row.Money)
    }

    return nil
}
