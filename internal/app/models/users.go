package models

// User
type User struct {
	Login          string `db:"login" json:"login"`
	HashedPassword string `db:"hashed_password" json:"hashed_password"`
	IsAdmin        bool   `db:"is_admin" json:"is_admin"`
}

// UserRequest
type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Validate
func (r *UserRequest) Validate() bool {
	return len(r.Login) >= 1 && len(r.Login) <= 50 && len(r.Password) >= 6 && len(r.Password) <= 50
}