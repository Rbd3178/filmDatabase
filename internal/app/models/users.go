package models

// User
type User struct {
	Login          string
	HashedPassword string
	Role        string
}

// UserRequest
type UserRequest struct {
	Login    string
	Password string
}
