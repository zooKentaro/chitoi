package core

import (
    redigo "github.com/garyburd/redigo/redis"
    "github.com/jmoiron/sqlx"
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/database"
    "github.com/uenoryo/chitoi/fluent"
    masterdata "github.com/uenoryo/chitoi/masterdata/repository"
    "github.com/uenoryo/chitoi/redis"
)

type Core struct {
    DB         *sqlx.DB
    Redis      redigo.Conn
    Logger     *fluent.Logger
    Masterdata *masterdata.Masterdata
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

    logger, err := fluent.Connect()
    if err != nil {
        return nil, errors.Wrap(err, "error connect fluent")
    }

    repositories := masterdata.New()
    if err := repositories.Load(dbConn); err != nil {
        return nil, errors.Wrap(err, "error load masterdata")
    }

    return &Core{
        DB:         dbConn,
        Redis:      redisConn,
        Logger:     logger,
        Masterdata: repositories,
    }, nil
}
