package row

import (
    "github.com/kayac/ddl-maker/dialect"
    "github.com/kayac/ddl-maker/dialect/mysql"
)

type UserRank struct {
    ID         uint32 `ddl:"auto" json:"id"`
    Rank       uint32 `json:"rank"`
    Assets     uint64 `json:"assets,string"`
    NormalRate uint32 `db:"normal_rate" json:"normal_rate"`
    HardRate   uint32 `db:"hard_rate" json:"hard_rate"`
}

func (r UserRank) PrimaryKey() dialect.PrimaryKey {
    return mysql.AddPrimaryKey("id")
}
