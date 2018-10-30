package row

import (
	"time"

	"github.com/kayac/ddl-maker/dialect"
	"github.com/kayac/ddl-maker/dialect/mysql"
)

type User struct {
	ID             uint64    `ddl:"auto" json:"id,string"`
	Name           string    `json:"name"`
	Token          string    `json:"token"`
	LastLoginAt    time.Time `db:"last_login_at" json:"last_login_at"`
	Rank           uint32    `ddl:"default=1" json:"rank"`
	Money          int64     `ddl:"default=0" json:"money"`
	Stamina        uint32    `ddl:"default=0" json:"stamina"`
	BestScore      uint64    `db:"best_score" ddl:"default=0" json:"best_score,string"`
	BestTotalScore uint64    `db:"best_total_score" ddl:"default=0" json:"best_total_score,string"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

func (u User) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}
