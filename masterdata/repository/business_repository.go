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
	_rows    []*row.Business
	_rowByID map[uint32]*row.Business
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
	repo._rowByID = make(map[uint32]*row.Business, len(rows))
	for _, row := range rows {
		repo._rowByID[row.ID] = row
	}
	return nil
}

// All ...
func (repo *BusinessRepository) All() []*row.Business {
	return repo._rows
}

// FindByID ...
func (repo *BusinessRepository) FindByID(id uint32) (*row.Business, error) {
	row, ok := repo._rowByID[id]
	if !ok {
		return nil, ErrNotFound
	}
	return row, nil
}
