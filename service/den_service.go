package service

import (
    "github.com/uenoryo/chitoi/core"
)

type DenService interface{}

type denService struct {
    Core *core.Core
}

// NewDenService is XXX
func NewDenService(core *core.Core) DenService {
    return &DenService{
        Core: core,
    }
}
