// Package model provides data Structures
package model

import "github.com/google/uuid"

// Auth struct represents the authentication mechanism
type Auth struct {
	Login    string `json:"login"`
	Password []byte `json:"password"`
}

// Login struct represents the authentication mechanism
type Login struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// NewUser struct represents the authentication mechanism
type NewUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// User struct represents the authentication mechanism
type User struct {
	Login    string `json:"login"`
	Password []byte `json:"password"`
	Username string `json:"username"`
}

// JWTResponse struct represents the JWT logic
type JWTResponse struct {
	ID           uuid.UUID `json:"id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken []byte    `json:"refresh_token"`
}
