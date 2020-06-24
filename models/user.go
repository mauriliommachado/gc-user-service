package models

// User model
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	ExpireAt string `json:"exp"`
}
