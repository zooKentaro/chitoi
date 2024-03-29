package model

import (
    "fmt"
    "time"

    "github.com/pkg/errors"
    uuid "github.com/satori/go.uuid"
    "github.com/uenoryo/chitoi/constant"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/database/row"
    "github.com/uenoryo/chitoi/lib/helper"
)

const (
    CreateUserSQL             = "INSERT INTO `user` (`name`, `token`, `last_login_at`, `rank`, `money`, `stamina`, `created_at`) VALUES (?,?,?,?,?,?,?)"
    LockForUpdateUserSQL      = "SELECT * FROM user WHERE id = ? FOR UPDATE"
    UpdateUserByLoginSQL      = "UPDATE user SET last_login_at = ?, money = ? WHERE id = ?"
    UpdateUserByFinishGameSQL = "UPDATE user SET stamina = ?, money = ? WHERE id = ?"
    InitialUserName           = ""
    SessionKeyPrefix          = "CHITOI-LOGIN-SESSION"
)

func CreateNewUser(core *core.Core) (*User, error) {
    v4, err := uuid.NewV4()
    if err != nil {
        return nil, errors.Wrap(err, "error uuid new v4")
    }
    token := v4.String()
    now := time.Now()

    userRow := &row.User{
        Token:       token,
        LastLoginAt: now,
        Money:       constant.DefaultMoney,
        Rank:        1,
        Stamina:     constant.DefaultStamina,
        CreatedAt:   now,
    }

    res, err := core.DB.Exec(CreateUserSQL, InitialUserName, token, now, userRow.Rank, userRow.Money, userRow.Stamina, now)
    if err != nil {
        return nil, errors.Wrap(err, "error create user")
    }

    id, err := res.LastInsertId()
    if err != nil {
        return nil, errors.Wrap(err, "error last insert id")
    }
    userRow.ID = uint64(id)

    return NewUser(core, userRow), nil
}

func (u *User) Login() (string, bool, error) {
    v4, err := uuid.NewV4()
    if err != nil {
        return "", false, errors.Wrap(err, "error uuid new v4")
    }
    sessionID := v4.String()
    key := fmt.Sprintf("%s:%s", SessionKeyPrefix, sessionID)
    if _, err := u.core.Redis.Do("SET", key, u.Row.ID); err != nil {
        return "", false, errors.Wrap(err, "error set session")
    }

    // 有効期限 2週間
    expire := 60 * 60 * 24 * 14
    if _, err := u.core.Redis.Do("EXPIRE", key, expire); err != nil {
        return "", false, errors.Wrap(err, "error set expire")
    }

    // 日付が変わっていたら収益を付与して last_login_at を更新
    if _, err := u.core.DB.Exec(LockForUpdateUserSQL, u.Row.ID); err != nil {
        return "", false, errors.Wrap(err, "error lock for update")
    }

    isTodayFirstLogin := helper.BeginningOfDayFromTime(time.Now()).After(u.Row.LastLoginAt)
    if isTodayFirstLogin {
        selected, err := NewUserBusinessRepository(u.core).SelectByUserID(u.Row.ID)
        if err != nil {
            return "", false, errors.Wrap(err, "error select user business by user id")
        }
        ubs := UserBusinesses(selected)
        profits, err := ubs.Profits()
        if err != nil {
            return "", false, errors.Wrap(err, "error profits")
        }

        u.Row.Money += profits

        if _, err := u.core.DB.Exec(UpdateUserByLoginSQL, time.Now(), u.Row.Money, u.Row.ID); err != nil {
            return "", false, errors.Wrap(err, "error update user data")
        }
    }

    return sessionID, isTodayFirstLogin, nil
}

func (u *User) UpdateRecord(bestScore, bestTotalScore uint64) error {
    if _, err := u.core.DB.Exec("SELECT * FROM user WHERE id = ? FOR UPDATE", u.Row.ID); err != nil {
        return errors.Wrap(err, "error lock for update")
    }

    if u.Row.BestScore > bestScore && u.Row.BestTotalScore > bestTotalScore {
        return nil
    }

    if u.Row.BestScore < bestScore {
        u.Row.BestScore = bestScore
    }

    if u.Row.BestTotalScore < bestTotalScore {
        u.Row.BestTotalScore = bestTotalScore
    }

    q := "UPDATE user SET best_score = ?, best_total_score = ? WHERE id = ?"
    if _, err := u.core.DB.Exec(q, u.Row.BestScore, u.Row.BestTotalScore, u.Row.ID); err != nil {
        return errors.Wrap(err, "error update user data")
    }
    return nil
}

// RankupMaybe は現在の資産に応じてランクを上げる
func (u *User) RankupMaybe() (bool, error) {
    assets, err := u.Assets()
    if err != nil {
        return false, errors.Wrap(err, "error assets")
    }

    rank := uint32(1)
    for _, ur := range u.core.Masterdata.UserRank.All() {
        if assets < ur.Assets {
            break
        }
        rank = ur.Rank
    }

    if rank <= u.Row.Rank {
        return false, nil
    }

    u.Row.Rank = rank

    if _, err := u.core.DB.Exec("UPDATE user SET rank = ? WHERE id = ?", u.Row.Rank, u.Row.ID); err != nil {
        return false, errors.Wrap(err, "error update user")
    }
    return true, nil
}

// exhaustStamina はスタミナを1つ消費する
func (u *User) exhaustStamina() error {
    if u.Row.Stamina == 0 {
        return errors.New("stamina is 0")
    }

    u.Row.Stamina--
    return nil
}

// Assets は User の総資産額を返す
func (u *User) Assets() (uint64, error) {
    userBusinesses, err := NewUserBusinessRepository(u.core).SelectByUserID(u.Row.ID)
    if err != nil {
        return 0, errors.Wrap(err, "error select user business by user id")
    }
    ubs := UserBusinesses(userBusinesses)

    businesses, err := ubs.Businesses()
    if err != nil {
        return 0, errors.Wrap(err, "error businesses")
    }
    businessByID := make(map[uint32]*Business, len(businesses))
    for _, b := range businesses {
        businessByID[b.Row.ID] = b
    }

    assets := uint64(0)
    for _, ub := range userBusinesses {
        price, err := businessByID[ub.Row.BusinessID].Price(ub.Row)
        if err != nil {
            return 0, errors.Wrap(err, "error price")
        }
        assets += price
    }
    return assets, nil
}

// getOrLoseMoney はお金を取得(消費)する
func (u *User) getOrLoseMoney(amount int64) {
    u.Row.Money += amount
}

// spendMoney はお金を消費する
func (u *User) spendMoney(amount uint64) {
    u.Row.Money -= int64(amount)
}
