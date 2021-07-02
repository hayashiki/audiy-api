package entity

import (
	"cloud.google.com/go/datastore"
	"time"
)

type Star struct {
	ID        string    `json:"id" datastore:"-"`
	UserID      *User     `json:"user_id"`
	//User      *User     `json:"user"`
	Audio     *datastore.Key `json:"audio_id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Star) IsNode() {}


