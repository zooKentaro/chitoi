package core

import (
    "database/sql"

    redigo "github.com/garyburd/redigo/redis"
    "github.com/jmoiron/sqlx"
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/database"
    "github.com/uenoryo/chitoi/database/row"
    "github.com/uenoryo/chitoi/redis"
)

type Core struct {
    DB         *sqlx.DB
    Redis      redigo.Conn
    Masterdata Masterdata
}

type Masterdata struct {
    Businesses []*row.Business
    UserRanks  []*row.UserRank
}

func New() (*Core, error) {
    dbConn, err := database.Connect()
    if err != nil {
        return nil, errors.Wrap(err, "error connect database")
    }

    redisConn, err := redis.Connect()
    if err != nil {
        return nil, errors.Wrap(err, "error connect redis")
    }

    return &Core{
        DB:    dbConn,
        Redis: redisConn,
    }, nil
}

func (core *Core) LoadMasterdata() error {
    businesses := []*row.Business{}
    if err := core.DB.Select(&businesses, "SELECT * FROM business"); err != nil {
        if err != sql.ErrNoRows {
            return errors.Wrap(err, "error select all business")
        }
    }

    userRanks := []*row.UserRank{}
    if err := core.DB.Select(&userRanks, "SELECT * FROM user_rank"); err != nil {
        if err != sql.ErrNoRows {
            return errors.Wrap(err, "error select all user_rank")
        }
    }

    core.Masterdata.Businesses = businesses
    core.Masterdata.UserRanks = userRanks
    return nil
}
