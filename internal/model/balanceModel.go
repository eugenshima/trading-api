package model

import "github.com/google/uuid"

type Balance struct {
	ID      uuid.UUID `json:"id"`
	Balance float64   `json:"balance"`
}
