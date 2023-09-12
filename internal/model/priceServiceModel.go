// Package model provides model structures
package model

// StreamedShares represents streams
type StreamedShares struct {
	Share []string `json:"share"`
}

// Shares represents shares kekw
type Shares struct {
	ShareName  string      `json:"share_name"`
	SharePrice interface{} `json:"price"`
}
