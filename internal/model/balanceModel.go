// Package model provides data Structures
package model

import "github.com/google/uuid"

// Balance struct represents the current balance
type Balance struct {
	ID      uuid.UUID `json:"id"`
	Balance float64   `json:"balance"`
}
