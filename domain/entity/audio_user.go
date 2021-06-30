package entity

import (
	"time"

	"cloud.google.com/go/datastore"
)

const AudioUserKind = "AudioUsers"

type AudioUser struct {
	Key       *datastore.Key `datastore:"__key__"`
	ID        int64          `json:"id" datastore:"-"`
	User      int64          `json:"user_id"`
	Audio     *datastore.Key `json:"audio_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (a AudioUser) GetKey() *datastore.Key {
	return datastore.IncompleteKey(AudioUserKind, nil)
}

func (a AudioUser) SetID(key *datastore.Key) {
	a.ID = key.ID
}

func NewAudioUser(userID int64, audioID string) *AudioUser {
	audioKey := GetAudioKey(audioID)
	au := &AudioUser{
		User:  userID,
		Audio: audioKey,
	}
	au.Key = au.GetKey()
	return au
}
