package data

import (
    "github.com/uenoryo/chitoi/database/row"
)

// BusinessListResponse is XXX
type BusinessListResponse struct {
    Businesses []*row.Business `json:"businesses"`
}

// BusinessBuyRequest is XXX
type BusinessBuyRequest struct {
    SessionID  string `json:"session_id"`
    BusinessID uint32 `json:"business_id"`
}

// BusinessBuyResponse is XXX
type BusinessBuyResponse struct{}
