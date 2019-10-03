package model

import (
    "database/sql"
    "fmt"

    "github.com/garyburd/redigo/redis"
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/database/row"
)

const (
    FindUserByIDSQL    = "SELECT * FROM user WHERE id = ?"
    FindUserByTokenSQL = "SELECT * FROM user WHERE token = ?"
)

type UserRepository struct {
    core *core.Core
}

func NewUserRepository(core *core.Core) *UserRepository {
    return &UserRepository{core: core}
}

func NewUser(core *core.Core, row *row.User) *User {
    user := &User{
        core: core,
        Row:  row,
    }
    user.Room = &UserRoomBehavior{core, user}

    return user
}

type User struct {
    Row  *row.User
    core *core.Core

    Room *UserRoomBehavior
}

func (repo *UserRepository) FindByToken(token string) (*User, error) {
    userRow := row.User{}
    err := repo.core.DB.Get(&userRow, FindUserByTokenSQL, token)
    switch {
    case err == sql.ErrNoRows:
        return nil, errors.Wrap(err, "user is not found")
    case err != nil:
        return nil, err
    }
    return NewUser(repo.core, &userRow), nil
}

func (repo *UserRepository) FindByID(id uint64) (*User, error) {
    userRow := row.User{}
    err := repo.core.DB.Get(&userRow, FindUserByIDSQL, id)
    switch {
    case err == sql.ErrNoRows:
        return nil, errors.Wrap(err, "user is not found")
    case err != nil:
        return nil, err
    }
    return NewUser(repo.core, &userRow), nil
}

func (repo *UserRepository) FindBySessionID(sessionID string) (*User, error) {
    key := fmt.Sprintf("%s:%s", SessionKeyPrefix, sessionID)
    userID, err := redis.Uint64(repo.core.Redis.Do("GET", key))
    if err != nil {
        return nil, errors.Wrap(err, "error get user id by session id")
    }

    return repo.FindByID(userID)
}
