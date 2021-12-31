package model

import (
	"time"
)

type User struct {
	ID        string         `json:"id"`
	Email     string         `json:"email"`
	Name      string         `json:"name"`
	PhotoURL  string         `json:"photoURL"`
	ProviderID string         `json:"provider_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (User) IsNode() {}

// TODO: add providerID
func NewUser(id string, email string, name string, photoURL string) *User {
	return &User{
		ID:        id,
		Email:     email,
		Name:      name,
		PhotoURL:  photoURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
