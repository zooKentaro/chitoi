package model

import (
    "database/sql"
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
    Row *row.User
}

func (repo *UserRepository) FindByToken(token string) (*User, error) {
    userRow := row.User{}
    if err := repo.core.DB.Get(&userRow, "SELECT * FROM user WHERE token = ?", token); err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.Wrap(err, "user is not found")
        }
        return nil, err
    }
    return &User{Row: &userRow}, nil
}
