package masterdata

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/uenoryo/chitoi/database/row"
)

const (
	selectAllBusinessSQL = "SELECT * FROM business"
)

// BusinessRepository (､´･ω･)▄︻┻┳═一
type BusinessRepository struct {
	_rows []*row.Business
}

// Load ...
func (repo *BusinessRepository) Load(db *sqlx.DB) error {
	rows := []*row.Business{}
	if err := db.Select(&rows, selectAllBusinessSQL); err != nil {
		if err != sql.ErrNoRows {
			return errors.Wrapf(err, "error select all business, sql:%s", selectAllBusinessSQL)
		}
	}
	repo._rows = rows
	return nil
}

// All ...
func (repo *BusinessRepository) All() []*row.Business {
	return repo._rows
}
