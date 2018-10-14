package service

import (
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/model"
)

type AuthService interface {
    Authenticate(sessionID string) (*model.User, error)
}

type authService struct {
    Core *core.Core
}

// NewAuthService is XXX
func NewAuthService(core *core.Core) AuthService {
    return &authService{
        Core: core,
    }
}

// Authenticate is XXX
func (s *authService) Authenticate(sessionID string) (*model.User, error) {
    user, err := model.NewUserRepository(s.Core).FindBySessionID(sessionID)
    if err != nil {
        return nil, errors.Wrap(err, "error find by session id")
    }
    return user, nil
}
