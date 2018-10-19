package model

import (
    "database/sql"

    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/database/row"
)

type UserBusinessRepository struct {
    core *core.Core
}

func NewUserBusinessRepository(core *core.Core) *UserBusinessRepository {
    return &UserBusinessRepository{core: core}
}

func (repo *UserBusinessRepository) SelectByUserID(id uint64) ([]*UserBusiness, error) {
    ubRows := []*row.UserBusiness{}
    if err := repo.core.DB.Select(&ubRows, "SELECT * FROM user_business WHERE user_id = ?", id); err != nil {
        if err != sql.ErrNoRows {
            return nil, err
        }
    }

    ubs := make([]*UserBusiness, len(ubRows))
    for i, row := range ubRows {
        ubs[i] = &UserBusiness{
            core: repo.core,
            Row:  row,
        }
    }
    return ubs, nil
}

type UserBusiness struct {
    Row  *row.UserBusiness
    core *core.Core
}

type UserBusinesses []*UserBusiness

func (ubs UserBusinesses) Businesses() ([]*Business, error) {
    if len(ubs) == 0 {
        return []*Business{}, nil
    }

    businessIDs := make([]uint32, len(ubs))
    for i, ub := range ubs {
        businessIDs[i] = ub.Row.BusinessID
    }

    bs, err := NewBusinessRepository(ubs[0].core).SelectByIDs(businessIDs)
    if err != nil {
        return nil, errors.Wrap(err, "error select by ids")
    }
    return bs, nil
}
