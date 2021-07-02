package entity

import "cloud.google.com/go/datastore"

const UserKind = "Users"

type User struct {
	Key         *datastore.Key `datastore:"__key__"`
	ID    int64 `json:"id" datastore:"-"`
	Email string `json:"email" datastore:"email"`
}

func (User) IsNode() {}

func NewUser(id int64, email string) *User {
	return &User{
		ID:    id,
		Email: email,
	}
}

func GetUserKey(id int64) *datastore.Key {
	//entity := User{ID: id}
	return datastore.IDKey(UserKind, id, nil)
}

//func (r *User) GetKey() *datastore.Key {
//	if r.ID == "" {
//		return nil
//	}
//	return datastore.NameKey(AudioKind, r.ID, nil)
//}
