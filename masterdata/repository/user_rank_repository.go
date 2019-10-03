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
	_rows []*row.UserRank
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
	return nil
}

// All ...
func (repo *UserRankRepository) All() []*row.UserRank {
	return repo._rows
}
