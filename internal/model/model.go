// Package model has struct of essence
package model

import "github.com/google/uuid"

// Symbol contains fields that describe the shares of companies
type Symbol struct {
	ID  uuid.UUID
	Bid float32
	Ask float32
}
