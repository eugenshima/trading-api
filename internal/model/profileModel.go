package model

type Auth struct {
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
