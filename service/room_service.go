package service

import (
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
    return nil, nil
}
