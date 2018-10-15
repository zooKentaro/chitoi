package service

import (
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/data"
    "github.com/uenoryo/chitoi/model"
)

type UserService interface {
    Signup(*data.UserSignupRequest) (*data.UserSignupResponse, error)
    Login(*data.UserLoginRequest) (*data.UserLoginResponse, error)
    Info(*data.UserInfoRequest) (*data.UserInfoResponse, error)
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
func (u *userService) Signup(req *data.UserSignupRequest) (*data.UserSignupResponse, error) {
    user, err := model.CreateNewUser(u.Core)
    if err != nil {
        return nil, errors.Wrap(err, "error create new user")
    }

    sessionID, err := user.Login()
    if err != nil {
        return nil, errors.Wrap(err, "error login")
    }

    return &data.UserSignupResponse{
        User:      user.Row,
        SessionID: sessionID,
    }, nil
}

// Login is XXX
func (u *userService) Login(req *data.UserLoginRequest) (*data.UserLoginResponse, error) {
    user, err := model.NewUserRepository(u.Core).FindByToken(req.Token)
    if err != nil {
        return nil, errors.Wrap(err, "find user by token")
    }

    sessionID, err := user.Login()
    if err != nil {
        return nil, errors.Wrap(err, "error login")
    }

    return &data.UserLoginResponse{
        User:       user.Row,
        SessionID:  sessionID,
        Businesses: u.Core.Masterdata.Businesses,
    }, nil
}

// Info is XXX
func (u *userService) Info(req *data.UserInfoRequest) (*data.UserInfoResponse, error) {
    user, err := NewAuthService(u.Core).Authenticate(req.SessionID)
    if err != nil {
        return nil, errors.Wrap(err, "error authenticate user")
    }

    return &data.UserInfoResponse{
        User: user.Row,
    }, nil
}
