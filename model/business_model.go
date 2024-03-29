package model

import (
    "database/sql"
    "fmt"
    "math/rand"
    "strconv"
    "strings"
    "time"

    "github.com/garyburd/redigo/redis"
    "github.com/jmoiron/sqlx"
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

func (repo *BusinessRepository) FindByID(id uint32) (*Business, error) {
    businessRow := row.Business{}
    if err := repo.core.DB.Get(&businessRow, "SELECT * FROM business WHERE id = ?", id); err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.Wrap(err, "business is not found")
        }
        return nil, err
    }
    return &Business{
        Row:  &businessRow,
        core: repo.core,
    }, nil
}

func (repo *BusinessRepository) SelectByIDs(ids []uint32) ([]*Business, error) {
    businessRows := []*row.Business{}
    q, vs, err := sqlx.In("SELECT * FROM business WHERE id IN (?)", ids)
    if err != nil {
        return nil, errors.Wrap(err, "error sqlx in")
    }
    if err := repo.core.DB.Select(&businessRows, q, vs...); err != nil {
        if err != sql.ErrNoRows {
            return nil, err
        }
    }

    bs := make([]*Business, len(businessRows))
    for i, row := range businessRows {
        bs[i] = &Business{
            core: repo.core,
            Row:  row,
        }
    }
    return bs, nil
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
    key := "CHITOI-BUSINESS-LIST"
    exists := true
    dateAndNum, err := redis.String(repo.core.Redis.Do("GET", key))
    if err != nil {
        if err == redis.ErrNil {
            exists = false
        } else {
            return 0, errors.Wrap(err, "error get user id by session id")
        }
    }

    isOld := false
    prefNum := uint32(0)
    nowStr := time.Now().Format("2006-01-02")
    if exists {
        sp := strings.Split(dateAndNum, "_")
        if len(sp) < 2 {
            return 0, errors.Errorf("error get invalid date and pref num data %s", dateAndNum)
        }
        dateStr := sp[0]
        n, err := strconv.Atoi(sp[1])
        if err != nil {
            return 0, errors.Wrap(err, "error parse int")
        }
        prefNum = uint32(n)
        if dateStr != nowStr {
            isOld = true
        }
    }

    if !exists || isOld {
        rand.Seed(time.Now().UnixNano())
        prefNum = uint32(rand.Intn(15) + 32) // TODO: 現状は32 ~ 47 まで。マスターが入ったら 1 ~ 47にする

        if _, err := repo.core.Redis.Do("SET", key, fmt.Sprintf("%s_%d", nowStr, prefNum)); err != nil {
            return 0, errors.Wrap(err, "error set datetime and prefecture num")
        }
    }
    return prefNum, nil
}

type Business struct {
    Row  *row.Business
    core *core.Core
}

func (b *Business) IsOpen() error {
    pref, err := NewBusinessRepository(b.core).todaysPrefNum()
    if err != nil {
        return errors.Wrap(err, "error get today's pref num")
    }
    if b.Row.Prefecture != pref {
        return errors.Errorf("error business id:%d is not open, current open prefecture:%d", b.Row.ID, pref)
    }
    return nil
}

func (b *Business) Price(ub *row.UserBusiness) (uint64, error) {
    if ub == nil {
        return b.Row.PriceBase, nil
    }
    switch ub.Level {
    case 0:
        return b.Row.PriceBase, nil
    case 1:
        return b.Row.PriceBase, nil
    case 2:
        return b.Row.PriceLevel2, nil
    case 3:
        return b.Row.PriceLevel3, nil
    }
    return 0, errors.Errorf("error invalid business level:%d", ub.Level)
}

func (b *Business) NextPrice(ub *row.UserBusiness) (uint64, error) {
    if ub == nil {
        return b.Row.PriceBase, nil
    }
    switch ub.Level {
    case 0:
        return b.Row.PriceBase, nil
    case 1:
        return b.Row.PriceLevel2, nil
    case 2:
        return b.Row.PriceLevel3, nil
    case 3:
        return 0, nil
    }
    return 0, errors.Errorf("error invalid business level:%d", ub.Level)
}

func (b *Business) Profit(ub *row.UserBusiness) (int64, error) {
    if ub == nil {
        return 0, nil
    }
    switch ub.Level {
    case 0:
        return int64(float32(b.Row.PriceBase) * float32(b.Row.ReturnRateBase) / 10 / 100), nil
    case 1:
        return int64(float32(b.Row.PriceBase) * float32(b.Row.ReturnRateBase) / 10 / 100), nil
    case 2:
        return int64(float32(b.Row.PriceLevel2) * float32(b.Row.ReturnRateLevel2) / 10 / 100), nil
    case 3:
        return int64(float32(b.Row.PriceLevel3) * float32(b.Row.ReturnRateLevel3) / 10 / 100), nil
    }
    return 0, errors.Errorf("error invalid business level:%d", ub.Level)
}
