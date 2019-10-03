package masterdata

import "github.com/jmoiron/sqlx"

// Masterdata (､´･ω･)▄︻┻┳═一
type Masterdata struct {
	Business *BusinessRepository
	UserRank *UserRankRepository
}

// New ...
func New() *Masterdata {
	return &Masterdata{
		&BusinessRepository{},
		&UserRankRepository{},
	}
}

// Load ...
func (md *Masterdata) Load(db *sqlx.DB) error {
	if err := md.Business.Load(db); err != nil {
		return err
	}

	if err := md.UserRank.Load(db); err != nil {
		return err
	}
	return nil
}
