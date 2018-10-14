package data

// GameFinishRequest is XXX
type GameFinishRequest struct {
    SessionID string `json:"session_id"`
    Money     uint64 `json:"money,string"`
}

func (req *GameFinishRequest) GetSessionID() string {
    return req.SessionID
}

// GameFinishResponse is XXX
type GameFinishResponse struct{}
