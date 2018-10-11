package service

import (
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/data"
    "github.com/uenoryo/chitoi/model"
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
    user, err := model.CreateNewUser(u.Core)
    if err != nil {
        return nil, errors.Wrap(err, "error create new user")
    }
    return &data.UserSignupResponse{
        User: user,
    }, nil
}
