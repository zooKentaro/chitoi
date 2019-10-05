package row

import (
	"time"

	"github.com/kayac/ddl-maker/dialect"
	"github.com/kayac/ddl-maker/dialect/mysql"
)

type Contact struct {
	ID        uint64    `ddl:"auto" json:"id,string"`
	UserID    uint64    `db:"user_id" ddl:"user_id" json:"user_id"`
	Body      string    `ddl:"body" json:"body"`
	Status    uint32    `ddl:"status" json:"status"`
	Note      string    `ddl:"note" json:"note"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt time.Time `ddl:"null" db:"deleted_at" json:"deleted_at"`
}

func (r *Contact) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}

func (r *Contact) Indexes() dialect.Indexes {
	return dialect.Indexes{
		mysql.AddIndex("user_id_idx", "user_id"),
	}
}
