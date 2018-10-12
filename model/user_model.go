package model

import (
    "database/sql"
    "fmt"
    "time"

    "github.com/pkg/errors"
    uuid "github.com/satori/go.uuid"
    "github.com/uenoryo/chitoi/constant"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/database/row"
)

func CreateNewUser(core *core.Core) (*row.User, error) {
    token := uuid.NewV4().String()
    now := time.Now()

    userRow := &row.User{
        Token:       token,
        LastLoginAt: now,
        Money:       constant.DefaultMoney,
        Stamina:     constant.DefaultStamina,
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    q := "INSERT INTO `user` (`name`, `token`, `last_login_at`, `money`, `stamina`, `created_at`, `updated_at`) VALUES (?,?,?,?,?,?,?)"
    res, err := core.DB.Exec(q, "", token, now, userRow.Money, userRow.Stamina, now, now)
    if err != nil {
        return nil, errors.Wrap(err, "error create user")
    }

    id, err := res.LastInsertId()
    if err != nil {
        return nil, errors.Wrap(err, "error last insert id")
    }
    userRow.ID = uint64(id)

    return userRow, nil
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

func (u *User) Login() (string, error) {
    sessionID := uuid.NewV4().String()
    key := fmt.Sprintf("CHITOI-LOGIN-SESSION:%s", sessionID)
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
