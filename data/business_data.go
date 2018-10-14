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
    BusinessID uint32 `json:"business_id"`
}
