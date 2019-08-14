package packet

import "github.com/uenoryo/chitoi/database/row"

// RequestPacket は各クライアントから送信される1回分のデータ
type RequestPacket struct {
	SessionID  string `json:"session_id"`
	ActionType uint32 `json:"action_type"`
	SenderID   uint64
	RoomCode   uint32
}

// BloadcastPacket は全体に送信するデータ
type BloadcastPacket struct {
	*RequestPacket
	Player1 *row.User `json:"player1"`
	Player2 *row.User `json:"player2"`
}
