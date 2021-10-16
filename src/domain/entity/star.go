package entity

import (
	"time"

	"cloud.google.com/go/datastore"
)

const StarKind = "Star"

type Star struct {
	Key       *datastore.Key `datastore:"__key__"`
	ID        int64          `json:"id" datastore:"-"`
	UserKey   *datastore.Key `json:"user_key" datastore:"user_key"`
	AudioKey  *datastore.Key `json:"audio_key" datastore:"audio_key"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

func (Star) IsNode() {}

func NewStar(userID string, audioID string) *Star {
	audioKey := GetAudioKey(audioID)
	userKey := GetUserKey(userID)
	au := &Star{
		UserKey:   userKey,
		AudioKey:  audioKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return au
}
