package service

import (
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/data"
    "github.com/uenoryo/chitoi/model"
)

type BusinessService interface {
    List(*data.BusinessListRequest) (*data.BusinessListResponse, error)
}

type businessService struct {
    Core *core.Core
}

// NewBusinessService is XXX
func NewBusinessService(core *core.Core) BusinessService {
    return &businessService{
        Core: core,
    }
}

// List is XXX
func (s *businessService) List(req *data.BusinessListRequest) (*data.BusinessListResponse, error) {
    businesses, err := model.NewBusinessRepository(s.Core).TodaysBusiness()
    if err != nil {
        return nil, errors.Wrap(err, "error today's business")
    }

    return &data.BusinessListResponse{
        Businesses: businesses,
    }, nil
}
