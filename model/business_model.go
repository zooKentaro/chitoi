package model

import (
    "database/sql"

    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/core"
    "github.com/uenoryo/chitoi/database/row"
)

type BusinessRepository struct {
    core *core.Core
}

func NewBusinessRepository(core *core.Core) *BusinessRepository {
    return &BusinessRepository{core: core}
}

func (repo *BusinessRepository) TodaysBusiness() ([]*row.Business, error) {
    businessRows := []*row.Business{}
    if err := repo.core.DB.Select(&businessRows, "SELECT * FROM business WHERE prefecture = ?", 43); err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.Wrap(err, "business is not found")
        }
        return nil, err
    }
    return businessRows, nil
}
