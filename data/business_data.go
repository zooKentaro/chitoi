package data

import (
    "github.com/uenoryo/chitoi/database/row"
)

// BusinessListRequest is XXX
type BusinessListRequest struct{}

// BusinessListResponse is XXX
type BusinessListResponse struct {
    Businesses []*row.Business
}
