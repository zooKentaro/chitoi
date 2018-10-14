package row

import (
    "github.com/kayac/ddl-maker/dialect"
    "github.com/kayac/ddl-maker/dialect/mysql"
)

type Business struct {
    ID               uint64 `json:"id"`
    Prefecture       uint32 `json:"prefecture"`
    Name             string `json:"name"`
    PriceBase        uint32 `db:"price_base" json:"price_base"`
    PriceLevel2      uint32 `db:"price_level2" json:"price_level2"`
    PriceLevel3      uint32 `db:"price_level3" json:"price_level3"`
    ReturnRateBase   uint32 `db:"return_rate_base" json:"return_rate_base"`
    ReturnRateLevel2 uint32 `db:"return_rate_level2" json:"return_rate_level2"`
    ReturnRateLevel3 uint32 `db:"return_rate_level3" json:"return_rate_level3"`
}

func (b Business) PrimaryKey() dialect.PrimaryKey {
    return mysql.AddPrimaryKey("id")
}
