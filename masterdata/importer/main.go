package main

import (
	"log"

	"github.com/uenoryo/chitoi/masterdata/importer/importer"
)

func main() {
	log.Println("start import")

	importer := importer.NewImporter("path/to/csv")

	if err := importer.Run(); err != nil {
		log.Println("error import. err: ", err.Error())
	}

	log.Println("finish import")
}
