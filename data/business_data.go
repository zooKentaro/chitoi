package data

import (
    "github.com/uenoryo/chitoi/database/row"
)

// BusinessListResponse is XXX
type BusinessListResponse struct {
    Businesses []*row.Business `json:"businesses"`
}
