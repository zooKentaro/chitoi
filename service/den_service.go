package service

import (
    "github.com/uenoryo/chitoi/core"
)

type DenService interface {
}

type denService struct {
    Core *core.Core
}

// NewDenService (､´･ω･)▄︻┻┳═一
func NewDenService(core *core.Core) DenService {
    return &denService{
        Core: core,
    }
}
