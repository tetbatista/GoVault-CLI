package domain

import "time"

type User struct {
	ID        string
	Username  string
	Password  string
	Salt      string
	CreatedAt time.Time
}

type UserRepository interface {
	Create(user *User) error
	FindByUsername(username string) (*User, error)
}
