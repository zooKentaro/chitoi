package model

import (
    "database/sql"
    "fmt"
    "time"

    "github.com/garyburd/redigo/redis"
    "github.com/pkg/errors"
    uuid "github.com/satori/go.uuid"
    "github.com/uenoryo/chitoi/constant"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/database/row"
    "github.com/uenoryo/chitoi/lib/helper"
)

var sessionKeyPrefix = "CHITOI-LOGIN-SESSION"

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
        Stamina:     constant.DefaultStamina,
        CreatedAt:   now,
    }

    q := "INSERT INTO `user` (`name`, `token`, `last_login_at`, `rank`, `money`, `stamina`, `created_at`) VALUES (?,?,?,?,?,?,?)"
    res, err := core.DB.Exec(q, "", token, now, 1, userRow.Money, userRow.Stamina, now)
    if err != nil {
        return nil, errors.Wrap(err, "error create user")
    }

    id, err := res.LastInsertId()
    if err != nil {
        return nil, errors.Wrap(err, "error last insert id")
    }
    userRow.ID = uint64(id)

    return &User{
        Row:  userRow,
        core: core,
    }, nil
}

type UserRepository struct {
    core *core.Core
}

func NewUserRepository(core *core.Core) *UserRepository {
    return &UserRepository{core: core}
}

type User struct {
    Row  *row.User
    core *core.Core
}

func (repo *UserRepository) FindByToken(token string) (*User, error) {
    userRow := row.User{}
    if err := repo.core.DB.Get(&userRow, "SELECT * FROM user WHERE token = ?", token); err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.Wrap(err, "user is not found")
        }
        return nil, err
    }
    return &User{
        Row:  &userRow,
        core: repo.core,
    }, nil
}

func (repo *UserRepository) FindByID(id uint64) (*User, error) {
    userRow := row.User{}
    if err := repo.core.DB.Get(&userRow, "SELECT * FROM user WHERE id = ?", id); err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.Wrap(err, "user is not found")
        }
        return nil, err
    }
    return &User{
        Row:  &userRow,
        core: repo.core,
    }, nil
}

func (repo *UserRepository) FindBySessionID(sessionID string) (*User, error) {
    key := fmt.Sprintf("%s:%s", sessionKeyPrefix, sessionID)
    userID, err := redis.Uint64(repo.core.Redis.Do("GET", key))
    if err != nil {
        return nil, errors.Wrap(err, "error get user id by session id")
    }

    return repo.FindByID(userID)
}

func (u *User) Login() (string, bool, error) {
    v4, err := uuid.NewV4()
    if err != nil {
        return "", false, errors.Wrap(err, "error uuid new v4")
    }
    sessionID := v4.String()
    key := fmt.Sprintf("%s:%s", sessionKeyPrefix, sessionID)
    if _, err := u.core.Redis.Do("SET", key, u.Row.ID); err != nil {
        return "", false, errors.Wrap(err, "error set session")
    }

    // 有効期限 2週間
    expire := 60 * 60 * 24 * 14
    if _, err := u.core.Redis.Do("EXPIRE", key, expire); err != nil {
        return "", false, errors.Wrap(err, "error set expire")
    }

    // 日付が変わっていたら収益を付与して last_login_at を更新
    if _, err := u.core.DB.Exec("SELECT * FROM user WHERE id = ? FOR UPDATE", u.Row.ID); err != nil {
        return "", false, errors.Wrap(err, "error lock for update")
    }

    isTodayFirstLogin := helper.BeginningOfDayFromTime(time.Now()).After(u.Row.LastLoginAt)
    if isTodayFirstLogin {
        q := "UPDATE user SET last_login_at = ? WHERE id = ?"
        if _, err := u.core.DB.Exec(q, time.Now(), u.Row.ID); err != nil {
            return "", false, errors.Wrap(err, "error update user data")
        }
    }

    return sessionID, isTodayFirstLogin, nil
}

// GameData は1ゲームのデータを扱う
type GameData struct {
    Money int64
}

// GameFinish は1ゲーム終了時の動作を行う
func (u *User) GameFinish(data *GameData) error {
    if err := u.exhaustStamina(); err != nil {
        return errors.Wrap(err, "error exhaust stamina")
    }

    u.getOrLoseMoney(data.Money)

    if _, err := u.core.DB.Exec("SELECT * FROM user WHERE id = ? FOR UPDATE", u.Row.ID); err != nil {
        return errors.Wrap(err, "error lock for update")
    }

    if _, err := u.core.DB.Exec("UPDATE user SET stamina = ?, money = ? WHERE id = ?", u.Row.Stamina, u.Row.Money, u.Row.ID); err != nil {
        return errors.Wrap(err, "error update user data")
    }

    return nil
}

func (u *User) BusinessList() ([]*row.UserBusiness, error) {
    ubRows := []*row.UserBusiness{}
    if err := u.core.DB.Select(&ubRows, "SELECT * FROM user_business WHERE user_id = ?", u.Row.ID); err != nil {
        if err != sql.ErrNoRows {
            return nil, errors.Wrap(err, "error find user business")
        }
    }

    return ubRows, nil
}

func (u *User) BusinessBuy(business *Business) error {
    if err := business.IsOpen(); err != nil {
        return errors.Wrap(err, "error business is open")
    }

    if _, err := u.core.DB.Exec("SELECT * FROM user WHERE id = ? FOR UPDATE", u.Row.ID); err != nil {
        return errors.Wrap(err, "error lock for update")
    }

    exists := true
    ubRow := &row.UserBusiness{}
    q := "SELECT * FROM user_business WHERE user_id = ? AND business_id = ?"
    if err := u.core.DB.Get(ubRow, q, u.Row.ID, business.Row.ID); err != nil {
        if err == sql.ErrNoRows {
            exists = false
        } else {
            return errors.Wrap(err, "error find user business")
        }
    }

    if err := u.canBuy(ubRow, business); err != nil {
        return errors.Wrap(err, "error can buy")
    }

    nextPrice, err := business.NextPrice(ubRow)
    if err != nil {
        return errors.Wrap(err, "error next price")
    }
    u.spendMoney(nextPrice)

    if _, err := u.core.DB.Exec("UPDATE user SET money = ? WHERE id = ?", u.Row.Money, u.Row.ID); err != nil {
        return errors.Wrap(err, "error update user data")
    }

    if exists {
        q := "UPDATE user_business SET level = ?, last_buy_at = ? WHERE user_id = ? AND business_id = ?"
        if _, err := u.core.DB.Exec(q, ubRow.Level+1, time.Now(), u.Row.ID, business.Row.ID); err != nil {
            return errors.Wrap(err, "error update user data")
        }
    } else {
        q := "INSERT INTO user_business (user_id, business_id, level, last_buy_at) VALUES (?,?,?,?)"
        _, err := u.core.DB.Exec(q, u.Row.ID, business.Row.ID, 1, time.Now())
        if err != nil {
            return errors.Wrap(err, "error create user")
        }
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
    for _, ur := range u.core.Masterdata.UserRanks {
        if assets < ur.Assets {
            break
        }
        rank = ur.Rank
    }

    if rank <= u.Row.Rank {
        return false, nil
    }

    if _, err := u.core.DB.Exec("UPDATE user SET rank = ? WHERE id = ?", rank, u.Row.ID); err != nil {
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

// canBuy は購入可能な情状態かどうか以下のチェックを行います
// * 今日はまだ1度も購入していないこと
// * まだ Business Level が最大になっていないこと
// * 購入に必要な資金を持っていること
func (u *User) canBuy(ub *row.UserBusiness, business *Business) error {
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

    if u.Row.Money < int64(nextPrice) {
        return errors.Errorf("error money is not enough, want %d but current %d", nextPrice, u.Row.Money)
    }

    return nil
}
