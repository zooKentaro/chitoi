package row

import (
	"time"

	"github.com/kayac/ddl-maker/dialect"
	"github.com/kayac/ddl-maker/dialect/mysql"
)

type User struct {
	ID          uint64    `ddl:"auto" json:"id,string"`
	Name        string    `json:"name"`
	Token       string    `json:"token"`
	LastLoginAt time.Time `json:"last_login_at"`
	Money       uint64    `ddl:"default=0" json:"money"`
	Stamina     uint32    `ddl:"default=0" json:"stamina"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u User) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}
