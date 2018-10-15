package service

import (
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/data"
    "github.com/uenoryo/chitoi/model"
)

type BusinessService interface {
    List() (*data.BusinessListResponse, error)
    Buy(*data.BusinessBuyRequest) (*data.BusinessBuyResponse, error)
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
func (s *businessService) List() (*data.BusinessListResponse, error) {
    businesses, err := model.NewBusinessRepository(s.Core).TodaysBusiness()
    if err != nil {
        return nil, errors.Wrap(err, "error today's business")
    }

    return &data.BusinessListResponse{
        Businesses: businesses,
    }, nil
}

// Buy is XXX
func (s *businessService) Buy(req *data.BusinessBuyRequest) (*data.BusinessBuyResponse, error) {
    user, err := NewAuthService(s.Core).Authenticate(req.SessionID)
    if err != nil {
        return nil, errors.Wrap(err, "error authenticate user")
    }

    business, err := model.NewBusinessRepository(s.Core).FindByID(req.BusinessID)
    if err != nil {
        return nil, errors.Wrapf(err, "error find business by id %s", req.BusinessID)
    }

    if err := user.BusinessBuy(business); err != nil {
        return nil, errors.Wrap(err, "error buy business")
    }

    return &data.BusinessBuyResponse{}, nil
}
