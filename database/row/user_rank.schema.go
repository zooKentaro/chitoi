package row

import (
    "github.com/kayac/ddl-maker/dialect"
    "github.com/kayac/ddl-maker/dialect/mysql"
)

type UserRank struct {
    ID              uint32 `ddl:"auto" json:"id"`
    Rank            uint32 `json:"rank"`
    Assets          uint64 `json:"assets,string"`
    NormalRate      uint32 `db:"normal_rate" json:"normal_rate"`
    HardRate        uint32 `db:"hard_rate" json:"hard_rate"`
    Com1NormalLevel uint32 `db:"com1_normal_level" json:"com1_normal_level"`
    Com2NormalLevel uint32 `db:"com2_normal_level" json:"com2_normal_level"`
    Com3NormalLevel uint32 `db:"com3_normal_level" json:"com3_normal_level"`
    Com1HardLevel   uint32 `db:"com1_hard_level" json:"com1_hard_level"`
    Com2HardLevel   uint32 `db:"com2_hard_level" json:"com2_hard_level"`
    Com3HardLevel   uint32 `db:"com3_hard_level" json:"com3_hard_level"`
}

func (r UserRank) PrimaryKey() dialect.PrimaryKey {
    return mysql.AddPrimaryKey("id")
}
