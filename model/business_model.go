package model

import (
    "database/sql"
    "fmt"
    "math/rand"
    "strconv"
    "strings"
    "time"

    "github.com/garyburd/redigo/redis"
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
        prefNum = uint32(rand.Intn(8) + 40) // TODO: 現状は40 ~ 47 まで。マスターが入ったら 1 ~ 47にする

        if _, err := repo.core.Redis.Do("SET", key, fmt.Sprintf("%s_%d", nowStr, prefNum)); err != nil {
            return 0, errors.Wrap(err, "error set datetime and prefecture num")
        }
    }

    return prefNum, nil
}