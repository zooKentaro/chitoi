package row

import (
	"time"

	"github.com/kayac/ddl-maker/dialect"
	"github.com/kayac/ddl-maker/dialect/mysql"
)

type BroadcastMessage struct {
	ID        uint64    `ddl:"auto" json:"id,string"`
	Title     string    `ddl:"title" json:"title"`
	Body      string    `ddl:"body" json:"body"`
	OpenAt    time.Time `ddl:"open_at" json:"open_at"`
	CloseAt   time.Time `ddl:"close_at" json:"close_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt time.Time `ddl:"null" db:"deleted_at" json:"deleted_at"`
}

func (r *BroadcastMessage) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}

func (r *BroadcastMessage) Indexes() dialect.Indexes {
	return dialect.Indexes{
		mysql.AddIndex("open_at_close_at_idx", "open_at", "close_at"),
	}
}
