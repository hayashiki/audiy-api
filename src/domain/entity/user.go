package entity

import (
	"time"

	"cloud.google.com/go/datastore"
)

const UserKind = "User"

type User struct {
	Key       *datastore.Key `datastore:"__key__"`
	ID        string         `json:"id" datastore:"-"`
	Email     string         `json:"email" validate:"email" datastore:"email"`
	Name      string         `json:"name" datastore:"name"`
	PhotoURL  string         `json:"photoURL" validate:"required" datastore:"photoURL"`
	CreatedAt time.Time      `json:"created_at" datastore:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" datastore:"updated_at"`
}

func (User) IsNode() {}

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

func GetUserKey(id string) *datastore.Key {
	//entity := User{ID: id}
	return datastore.NameKey(UserKind, id, nil)
}
