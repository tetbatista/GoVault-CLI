package domain

import "time"

type Credential struct {
	ID        string
	UserID    string
	Site      string
	Username  string
	Password  string
	CreatedAt time.Time
}

type CredentialRepository interface {
	Create(credential *Credential) error
	FindAllByUserID(userID string) ([]Credential, error)
	FindByUserIDAndSite(userID string, site string) (*Credential, error)
	DeleteByUserIDAndSite(userID string, site string) error
}
