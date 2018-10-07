package service

import (
    "time"

    "github.com/pkg/errors"
    uuid "github.com/satori/go.uuid"
    "github.com/uenoryo/chitoi/constant"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/data"
)

type UserService interface {
    Signup(*data.UserSignupRequest) (*data.UserSignupResponse, error)
}

type userService struct {
    Core *core.Core
}

// NewUserService is XXX
func NewUserService(core *core.Core) UserService {
    return &userService{
        Core: core,
    }
}

// Signup is XXX
func (u *userService) Signup(*data.UserSignupRequest) (*data.UserSignupResponse, error) {
    now := time.Now()
    q := "INSERT INTO `user` (`name`, `token`, `last_login_at`, `money`, `stamina`, `created_at`, `updated_at`) VALUES (?,?,?,?,?,?,?)"
    if _, err := u.Core.DB.Exec(q, "", uuid.NewV4().String(), now, constant.DefaultMoney, constant.DefaultStamina, now, now); err != nil {
        return nil, errors.Wrap(err, "error create user")
    }
    return &data.UserSignupResponse{}, nil
}
