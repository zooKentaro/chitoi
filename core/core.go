package core

import (
    redigo "github.com/garyburd/redigo/redis"
    "github.com/jmoiron/sqlx"
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/database"
    "github.com/uenoryo/chitoi/redis"
)

type Core struct {
    DB    *sqlx.DB
    Redis redigo.Conn
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
