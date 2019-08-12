package model

import (
	"github.com/pkg/errors"
	"github.com/uenoryo/chitoi/core"
	"github.com/uenoryo/chitoi/database/row"
)

const (
	FindUserByIDsSQL = "SELECT * FROM user WHERE id IN (?)"
)

// NewUserRepository (､´･ω･)▄︻┻┳═一
func NewUserRepository(core *core.Core) *UserRepository {
	return &UserRepository{core}
}

// UserRepository (､´･ω･)▄︻┻┳═一
type UserRepository struct {
	core *core.Core
}

// FindByIDs は ids の user を取得する
func (repo *UserRepository) FindByIDs(ids []uint64) ([]*User, error) {
	if len(ids) == 0 {
		return []*User{}, nil
	}

	rows := make([]*row.User, 0, len(ids))
	if err := repo.core.DB.Get(&rows, FindUserByIDsSQL, ids); err != nil {
		return nil, errors.Wrapf(err, "error find user by id: %v, sql: %s", ids, FindUserByIDsSQL)
	}

	users := make([]*User, len(rows))
	for i, row := range rows {
		users[i] = NewUser(repo.core, row)
	}
	return users, nil
}

// User (､´･ω･)▄︻┻┳═一
type User struct {
	core *core.Core
	Row  *row.User
}

// NewUser (､´･ω･)▄︻┻┳═一
func NewUser(core *core.Core, row *row.User) *User {
	return &User{
		core,
		row,
	}
}
