package packet

// Card はトランプのカード
type Card struct {
	ID     uint32
	Mark   uint32
	Number uint32
}

// Deck は山札
type Deck struct {
	Cards []*Card
}
