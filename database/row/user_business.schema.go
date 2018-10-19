package row

import (
    "time"

    "github.com/kayac/ddl-maker/dialect"
    "github.com/kayac/ddl-maker/dialect/mysql"
)

type UserBusiness struct {
    ID         uint64    `ddl:"auto" json:"id,string"`
    UserID     uint64    `db:"user_id" json:"user_id"`
    BusinessID uint32    `db:"business_id" json:"business_id"`
    Level      uint32    `db:"level" json:"level"`
    LastBuyAt  time.Time `db:"last_buy_at" json:"last_buy_at"`
}

func (r UserBusiness) PrimaryKey() dialect.PrimaryKey {
    return mysql.AddPrimaryKey("id")
}

// Indexes is XXX
func (r *UserBusiness) Indexes() dialect.Indexes {
    return dialect.Indexes{
        mysql.AddUniqueIndex("user_business_idx", "user_id", "business_id"),
    }
}
