package main

import (
    "io/ioutil"
    "log"
    "strings"

    _ "github.com/go-sql-driver/mysql"
    "github.com/pkg/errors"
    "github.com/uenoryo/chitoi/database"
    "github.com/uenoryo/chitoi/env"
    "github.com/uenoryo/hamster"
)

const (
    masterdataDir = "./masterdata/csv/"
)

func main() {
    if err := env.Load(); err != nil {
        log.Fatal("error load env, error: ", err.Error())
        return
    }

    db, err := database.ConnectStandard()
    if err != nil {
        log.Fatalf("error connect database, err: %s", err.Error())
    }
    defer db.Close()

    files, err := csvFiles()
    if err != nil {
        log.Fatalf("error load csv files")
    }

    feed := make([]*hamster.Food, len(files))
    for i, filename := range files {
        sp := strings.Split(filename, ".")
        if len(sp) == 0 {
            log.Fatalf("error invalid filename %s", filename)
            return
        }
        tableName := sp[0]

        feed[i] = &hamster.Food{
            Table:    tableName,
            Filepath: masterdataDir + filename,
        }
    }
    ham := hamster.New(db, &hamster.Option{})

    if err := ham.Stuff(feed); err != nil {
        log.Fatal(err.Error())
    }
}

func csvFiles() ([]string, error) {
    files, err := ioutil.ReadDir(masterdataDir)
    if err != nil {
        return nil, errors.Wrap(err, "error read directory")
    }

    var res []string
    for _, f := range files {
        res = append(res, f.Name())
    }
    return res, nil
}
