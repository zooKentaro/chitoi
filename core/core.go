package core

import (
    "github.com/jmoiron/sqlx"
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/database"
)

type Core struct {
    DB *sqlx.DB
}

func New() (*Core, error) {
    conn, err := database.Connect()
    if err != nil {
        return nil, errors.Wrap(err, "error connect database")
    }
    return &Core{
        DB: conn,
    }, nil
}
