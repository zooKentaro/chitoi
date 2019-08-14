package packet

import "github.com/uenoryo/chitoi/database/row"

const (
	// ActionTypeUndefined ...
	ActionTypeUndefined = iota

	// ActionTypeSetupGame ...
	ActionTypeSetupGame
)

// RequestPacket は各クライアントから送信される1回分のデータ
type RequestPacket struct {
	SessionID  string     `json:"session_id"`
	ActionType ActionType `json:"action_type"`
	SenderID   uint64
	RoomCode   uint32
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

// ActionType ...
type ActionType int

// IsUndefined ...
func (at ActionType) IsUndefined() bool {
	return int(at) == ActionTypeUndefined
}

// IsEntryAsHost ...
func (at ActionType) IsEntryAsHost() bool {
	return int(at) == ActionTypeEntryAsHost
}

// IsEntryAsGuest ...
func (at ActionType) IsEntryAsGuest() bool {
	return int(at) == ActionTypeEntryAsGuest
}
