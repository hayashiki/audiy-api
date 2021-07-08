package entity

import (
	"cloud.google.com/go/datastore"
	"time"
)

type Star struct {
	ID        string    `json:"id" datastore:"-"`
	UserKey    *datastore.Key `json:"user_key" datastore:"user_key"`
	AudioKey   *datastore.Key `json:"audio_key" datastore:"audio_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Star) IsNode() {}

func NewStar(userID int64, audioID string) *Star {
	audioKey := GetAudioKey(audioID)
	userKey := GetUserKey(userID)
	au := &Star{
		UserKey:  userKey,
		AudioKey: audioKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return au
}
