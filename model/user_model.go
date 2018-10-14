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
)

var sessionKeyPrefix = "CHITOI-LOGIN-SESSION"

func CreateNewUser(core *core.Core) (*User, error) {
    token := uuid.NewV4().String()
    now := time.Now()

    userRow := &row.User{
        Token:       token,
        LastLoginAt: now,
        Money:       constant.DefaultMoney,
        Stamina:     constant.DefaultStamina,
        CreatedAt:   now,
    }

    q := "INSERT INTO `user` (`name`, `token`, `last_login_at`, `money`, `stamina`, `created_at`) VALUES (?,?,?,?,?,?,?)"
    res, err := core.DB.Exec(q, "", token, now, userRow.Money, userRow.Stamina, now)
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

func (u *User) Login() (string, error) {
    sessionID := uuid.NewV4().String()
    key := fmt.Sprintf("%s:%s", sessionKeyPrefix, sessionID)
    if _, err := u.core.Redis.Do("SET", key, u.Row.ID); err != nil {
        return "", errors.Wrap(err, "error set session")
    }

    // 有効期限 2週間
    expire := 60 * 60 * 24 * 14
    if _, err := u.core.Redis.Do("EXPIRE", key, expire); err != nil {
        return "", errors.Wrap(err, "error set expire")
    }

    return sessionID, nil
}

// GameData は1ゲームのデータを扱う
type GameData struct {
    Money uint64
}

// GameFinish は1ゲーム終了時の動作を行う
func (u *User) GameFinish(data *GameData) error {
    if err := u.exhaustStamina(); err != nil {
        return errors.Wrap(err, "error exhaust stamina")
    }

    u.getMoney(data.Money)

    if _, err := u.core.DB.Exec("SELECT * FROM user WHERE id = ? FOR UPDATE", u.Row.ID); err != nil {
        return errors.Wrap(err, "error lock for update")
    }

    if _, err := u.core.DB.Exec("UPDATE user SET stamina = ?, money = ? WHERE id = ?", u.Row.Stamina, u.Row.Money, u.Row.ID); err != nil {
        return errors.Wrap(err, "error update user data")
    }

    return nil
}

// exhaustStamina はスタミナを1つ消費する
func (u *User) exhaustStamina() error {
    if u.Row.Stamina == 0 {
        return errors.New("stamina is 0")
    }

    u.Row.Stamina--
    return nil
}

// getMoney はお金を取得する
func (u *User) getMoney(amount uint64) {
    u.Row.Money += amount
}
