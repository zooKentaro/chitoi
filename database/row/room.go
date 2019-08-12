package row

// PlayerIDs は room の Player1 ~ Player4 のIDを返す
func (r *Room) PlayerIDs() []uint64 {
    res := make([]uint64, 0, 4)
    for _, id := range []uint64{r.Player1ID, r.Player2ID, r.Player3ID, r.Player4ID} {
        if id == 0 {
            continue
        }
        res = append(res, id)
    }
    return res
}
