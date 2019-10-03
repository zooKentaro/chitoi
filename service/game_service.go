package service

import (
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/data"
    "github.com/uenoryo/chitoi/model"
)

type GameService interface {
    Finish(*data.GameFinishRequest) (*data.GameFinishResponse, error)
}

type gameService struct {
    Core *core.Core
}

// NewGameService is XXX
func NewGameService(core *core.Core) GameService {
    return &gameService{
        Core: core,
    }
}

// Finish is XXX
func (s *gameService) Finish(req *data.GameFinishRequest) (*data.GameFinishResponse, error) {
    user, err := NewAuthService(s.Core).Authenticate(req.GetSessionID())
    if err != nil {
        return nil, errors.Wrap(err, "error authenticate user")
    }

    gd := &model.GameData{
        Money: req.Money,
    }

    if err := user.Game.Finish(gd); err != nil {
        return nil, errors.Wrapf(err, "error game finish, game data %+v", gd)
    }

    return &data.GameFinishResponse{}, nil
}
