package packet

// Card はトランプのカード
type Card struct {
	ID     uint32 `json:"id"`
	Mark   uint32 `json:"mark"`
	Number uint32 `json:"number"`
}

// Deck は山札
type Deck struct {
	Cards []*Card `json:"cards"`
}

type TurnTable struct {
	PlayerIDs []uint32 `json:"player_ids"`
}
