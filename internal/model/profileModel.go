package model

import "github.com/google/uuid"

type Auth struct {
	Login    string `json:"login"`
	Password []byte `json:"password"`
}

type Login struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type NewUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type User struct {
	Login    string `json:"login"`
	Password []byte `json:"password"`
	Username string `json:"username"`
}

type JWTResponse struct {
	ID           uuid.UUID `json:"id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}
