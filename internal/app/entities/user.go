package entities

// User
type User struct {
	Login          string `db:"login"`
	HashedPassword string `db:"hashed_password"`
	IsAdmin        string `db:"is_admin"`
}
