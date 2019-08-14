package packet

import "github.com/uenoryo/chitoi/database/row"

const (
	// MethodUndefined ...
	MethodUndefined = iota

	// MethodSetupGame ...
	MethodSetupGame
)

// RequestPacket は各クライアントから送信される1回分のデータ
type RequestPacket struct {
	SessionID string `json:"session_id"`
	Method    Method `json:"method"`
	SenderID  uint64
	RoomCode  uint32
	*SetupGameRequestPacket
}

// BloadcastPacket は全体に送信するデータ
type BloadcastPacket struct {
	*RequestPacket
	Player1 *row.User `json:"player1"`
	Player2 *row.User `json:"player2"`
}

type SetupGameRequestPacket struct {
	Deck      []*Deck
	TurnTable TurnTable
}

// Method ...
type Method int

// IsUndefined ...
func (m Method) IsUndefined() bool {
	return int(m) == MethodUndefined
}

// IsSetupGame ...
func (m Method) IsSetupGame() bool {
	return int(m) == MethodSetupGame
}
