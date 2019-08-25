package packet

import "github.com/uenoryo/chitoi/database/row"

const (
	// MethodUndefined ...
	MethodUndefined Method = iota
	// MethodSetupGame ゲームの初期化が完了し、ゲームを開始する
	MethodSetupGame
	// EntryPlayer Playerの新規接続
	MethodEntryPlayer
	// ExitPlayer Playerの接続解除
	MethodExitPlayer
	// Put カードを場に出した
	MethodPut
	// Draw カードを引いた
	MethodDraw
	// RefuseAttach カードを付け加えた
	MethodRefuseAttach
	// ChangeMark マークを変えた
	MethodChangeMark
)

// RequestPacket は各クライアントから送信される1回分のデータ
type RequestPacket struct {
	SessionID string `json:"session_id"`
	Method    Method `json:"method"`
	SenderID  uint64
	RoomCode  uint32
	*SetupGameRequestPacket
	*GameActionRequestPacket
}

// BloadcastPacket は全体に送信するデータ
type BloadcastPacket struct {
	*RequestPacket
	Player1 *row.User `json:"player1"`
	Player2 *row.User `json:"player2"`
}

// SetupGameRequestPacket はゲームをセットアップした時に乗せるデータ
type SetupGameRequestPacket struct {
	Deck      []*Deck   `json:"deck"`
	TurnTable TurnTable `json:"turn_table"`
}

// GameActionRequestPacket はゲーム内のプレイヤーの1回の行動のデータ
type GameActionRequestPacket struct {
	ActionType   uint32 `json:"action_type"`
	Mark         uint32 `json:"mark"`
	PutCardIndex uint32 `json:"put_card_index"`
}

// Method ...
type Method int

// IsUndefined ...
func (m Method) IsUndefined() bool {
	return m == MethodUndefined
}

// IsSetupGame ...
func (m Method) IsSetupGame() bool {
	return m == MethodSetupGame
}
