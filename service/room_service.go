package service

import (
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/data"
)

type RoomService interface {
    Create(*data.RoomCreateRequest) (*data.RoomCreateResponse, error)
}

type roomService struct {
    Core *core.Core
}

// NewRoomService is XXX
func NewRoomService(core *core.Core) RoomService {
    return &roomService{
        Core: core,
    }
}

// Create is XXX
func (s *roomService) Create(req *data.RoomCreateRequest) (*data.RoomCreateResponse, error) {
    user, err := NewAuthService(s.Core).Authenticate(req.SessionID)
    if err != nil {
        return nil, errors.Wrap(err, "error authenticate user")
    }

    if err := user.Room.Clean(); err != nil {
        return nil, errors.Wrap(err, "clean user room failed")
    }

    room, err := user.Room.Create()
    if err != nil {
        return nil, errors.Wrap(err, "create user room failed")
    }
    return &data.RoomCreateResponse{
        Room: room.Row,
    }, nil
}
