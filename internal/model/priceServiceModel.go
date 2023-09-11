package model

import "github.com/google/uuid"

type StreamedShares struct {
	ID     uuid.UUID `json:"id"`
	Shares []string  `json:"shares"`
}

type Shares struct {
	ShareName  string      `json:"share_name"`
	SharePrice interface{} `json:"price"`
}
