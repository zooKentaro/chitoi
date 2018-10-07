package row

import (
	"time"

	"github.com/kayac/ddl-maker/dialect"
	"github.com/kayac/ddl-maker/dialect/mysql"
)

type User struct {
	Id          uint64
	Name        string
	Token       string
	LastLoginAt time.Time
	money       uint64
	stamina     uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u User) PrimaryKey() dialect.PrimaryKey {
	return mysql.AddPrimaryKey("id")
}
