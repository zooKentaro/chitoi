package masterdata

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/uenoryo/chitoi/database/row"
)

const (
	selectAllUserRankSQL = "SELECT * FROM user_rank"
)

// UserRankRepository (､´･ω･)▄︻┻┳═一
type UserRankRepository struct {
	_rows    []*row.UserRank
	_rowByID map[uint32]*row.UserRank
}

// Load ...
func (repo *UserRankRepository) Load(db *sqlx.DB) error {
	rows := []*row.UserRank{}
	if err := db.Select(&rows, selectAllUserRankSQL); err != nil {
		if err != sql.ErrNoRows {
			return errors.Wrapf(err, "error select all user_rank, sql:%s", selectAllUserRankSQL)
		}
	}

	repo._rows = rows
	repo._rowByID = make(map[uint32]*row.UserRank, len(rows))
	for _, row := range rows {
		repo._rowByID[row.ID] = row
	}
	return nil
}

// All ...
func (repo *UserRankRepository) All() []*row.UserRank {
	return repo._rows
}

// FindByID ...
func (repo *UserRankRepository) FindByID(id uint32) (*row.UserRank, error) {
	row, ok := repo._rowByID[id]
	if !ok {
		return nil, ErrNotFound
	}
	return row, nil
}
