package model

import (
    "database/sql"
    "math/rand"
    "time"

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
    prefNum, err := repo.todaysPrefNum()
    if err != nil {
        return nil, errors.Wrap(err, "error today's pref num")
    }

    businessRows := []*row.Business{}
    if err := repo.core.DB.Select(&businessRows, "SELECT * FROM business WHERE prefecture = ?", prefNum); err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.Wrap(err, "business is not found")
        }
        return nil, err
    }
    return businessRows, nil
}

// todaysPrefNum は今日の都道府県番号を返す (日替わりのランダム)
func (repo *BusinessRepository) todaysPrefNum() (uint32, error) {
    rand.Seed(time.Now().UnixNano())
    num := uint32(rand.Intn(8) + 40) // TODO: 現状は40 ~ 47 まで。マスターが入ったら 1 ~ 47にする

    return num, nil
}
