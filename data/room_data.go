package data

import "github.com/uenoryo/chitoi/database/row"

// RoomCreateRequest is XXX
type RoomCreateRequest struct {
    SessionID string `json:"session_id"`
}

// RoomCreateResponse is XXX
type RoomCreateResponse struct {
    Room *row.Room
}
