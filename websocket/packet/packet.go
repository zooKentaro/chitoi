package packet

import "github.com/uenoryo/chitoi/database/row"

// RequestPacket は各クライアントから送信される1回分のデータ
type RequestPacket struct {
	SessionID    string `json:"session_id"`
	Method       Method `json:"method"`
	SenderID     uint64
	RoomCode     uint32
	TurnTable    *TurnTable `json:"turn_table"`
	RandomSeed   uint64     `json:"random_seed"`
	ActionType   ActionType `json:"action_type"`
	Mark         uint32     `json:"change_mark"`
	PutHandIndex int32      `json:"put_hand_index"`
}

// BroadcastPacket は全体に送信するデータ
type BroadcastPacket struct {
	*RequestPacket
	Player1 *row.User `json:"player1"`
	Player2 *row.User `json:"player2"`
}

// Method ...
type Method int

const (
	// MethodUndefined ...
	MethodUndefined Method = iota
	// MethodSetupGame ゲームの初期化が完了し、ゲームを開始する
	MethodSetupGame
	// EntryPlayer Playerの新規接続
	MethodEntryPlayer
	// ExitPlayer Playerの接続解除
	MethodExitPlayer
	// InGameAction インゲーム内の行動
	InGameAction
)

// IsUndefined ...
func (m Method) IsUndefined() bool {
	return m == MethodUndefined
}

// IsSetupGame ...
func (m Method) IsSetupGame() bool {
	return m == MethodSetupGame
}

// ActionType はインゲーム内の行動を識別する型
type ActionType int

const (
	// ActionTypePut カードを場に出した
	ActionTypePut ActionType = iota
	// ActionTypeDraw カードを引いた
	ActionTypeDraw
	// ActionTypeRefuseAttach カードを付け加えた
	ActionTypeRefuseAttach
	// ActionTypeChangeMark マークを変えた
	ActionTypeChangeMark
	// TurnAndPut ...
	TurnAndPut
	// Den かけようとした
	Den
)
