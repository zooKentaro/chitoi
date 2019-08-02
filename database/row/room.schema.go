package row

import (
    "time"

    "github.com/kayac/ddl-maker/dialect"
    "github.com/kayac/ddl-maker/dialect/mysql"
)

type Room struct {
    ID        uint64    `json:"id,string"`
    OwnerID   uint64    `db:"owner_id" json:"owner_id,string"`
    User1ID   uint64    `db:"user1_id" json:"user1_id,string"`
    User2ID   uint64    `db:"user2_id" json:"user2_id,string"`
    User3ID   uint64    `db:"user3_id" json:"user3_id,string"`
    User4ID   uint64    `db:"user4_id" json:"user4_id,string"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (r *Room) PrimaryKey() dialect.PrimaryKey {
    return mysql.AddPrimaryKey("id")
}

// Indexes is XXX
func (r *Room) Indexes() dialect.Indexes {
    return dialect.Indexes{
        mysql.AddIndex("owner_id_idx", "owner_id"),
    }
}
