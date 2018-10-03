package main

import (
	"fmt"
	"log"

	ddlmaker "github.com/kayac/ddl-maker"
	"github.com/pkg/errors"
	"github.com/uenoryo/chitoi/config"
	"github.com/uenoryo/chitoi/database/row"
)

const (
	driver  = "mysql"
	engine  = "InnoDB"
	charset = "utf8mb4"
)

func main() {
	if err := create(fmt.Sprintf("%s/database/sql/%s", config.Home(), "main.sql"), row.MainTableStructs()); err != nil {
		log.Fatal("error create main schema", err.Error())
	}
}

func create(schemaFile string, structs []interface{}) error {
	conf := ddlmaker.Config{
		DB: ddlmaker.DBConfig{
			Driver:  driver,
			Engine:  engine,
			Charset: charset,
		},
		OutFilePath: schemaFile,
	}

	ddlmaker, err := ddlmaker.New(conf)
	if err != nil {
		return errors.Wrap(err, "error new ddlmaker")
	}

	ddlmaker.AddStruct(structs...)

	if err = ddlmaker.Generate(); err != nil {
		return errors.Wrap(err, "error ddlmaker generate")
	}

	return nil
}
