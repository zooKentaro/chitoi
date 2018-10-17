package database

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Connect() (*sqlx.DB, error) {
	return sqlx.Connect(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo",
			os.Getenv("CHITOI_DB_USER"),
			os.Getenv("CHITOI_DB_PASS"),
			os.Getenv("CHITOI_DB_HOST"),
			os.Getenv("CHITOI_DB_PORT"),
			os.Getenv("CHITOI_DB_NAME"),
		),
	)
}
