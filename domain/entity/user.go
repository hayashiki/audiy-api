package entity

import (
	"time"

	"cloud.google.com/go/datastore"
)

const UserKind = "User"

type User struct {
	Key       *datastore.Key `datastore:"__key__"`
	ID        int64          `json:"id" datastore:"-"`
	Email     string         `json:"email" datastore:"email"`
	CreatedAt time.Time      `json:"created_at" datastore:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" datastore:"updated_at"`
}

func (User) IsNode() {}

func NewUser(id int64, email string) *User {
	return &User{
		ID:        id,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func GetUserKey(id int64) *datastore.Key {
	//entity := User{ID: id}
	return datastore.IDKey(UserKind, id, nil)
}
