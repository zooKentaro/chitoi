package row

import (
    "time"

    "github.com/kayac/ddl-maker/dialect"
    "github.com/kayac/ddl-maker/dialect/mysql"
)

type Room struct {
    ID        uint64    `json:"id,string"`
    Code      uint32    `json:"code"`
    OwnerID   uint64    `db:"owner_id" json:"owner_id,string"`
    Player1ID uint64    `db:"player1_id" json:"player1_id,string"`
    Player2ID uint64    `db:"player2_id" json:"player2_id,string"`
    Player3ID uint64    `db:"player3_id" json:"player3_id,string"`
    Player4ID uint64    `db:"player4_id" json:"player4_id,string"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    ExpiredAt time.Time `db:"expired_at" json:"expired_at"`
}

func (r *Room) PrimaryKey() dialect.PrimaryKey {
    return mysql.AddPrimaryKey("id")
}

// Indexes is XXX
func (r *Room) Indexes() dialect.Indexes {
    return dialect.Indexes{
        mysql.AddIndex("code_idx", "code"),
        mysql.AddIndex("owner_id_idx", "owner_id"),
    }
}
