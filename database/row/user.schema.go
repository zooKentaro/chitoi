package row

import (
	"time"

	"github.com/kayac/ddl-maker/dialect"
	"github.com/kayac/ddl-maker/dialect/mysql"
)

type User struct {
	Id          uint64 `ddl:"auto"`
	Name        string
	Token       string
	LastLoginAt time.Time
	money       uint64 `ddl:"default=0"`
	stamina     uint32 `ddl:"default=0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u User) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}
