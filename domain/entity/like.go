package entity

import (
	"cloud.google.com/go/datastore"
	"time"
)

const LikeKind = "Like"

type Like struct {
	Key       *datastore.Key `datastore:"__key__"`
	ID        int64    `json:"id" datastore:"-"`
	UserKey    *datastore.Key `json:"user_key" datastore:"user_key"`
	AudioKey   *datastore.Key `json:"audio_key" datastore:"audio_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Like) IsNode() {}

func NewLike(userID int64, audioID string) *Like {
	audioKey := GetAudioKey(audioID)
	userKey := GetUserKey(userID)
	au := &Like{
		UserKey:  userKey,
		AudioKey: audioKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return au
}
