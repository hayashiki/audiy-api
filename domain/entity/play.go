package entity

import (
	"time"

	"cloud.google.com/go/datastore"
)

const AudioUserKind = "Play"

//type Play struct {
//	ID        string    `json:"id"`
//	User      *User     `json:"user"`
//	Audio     *Audio    `json:"audio"`
//	CreatedAt time.Time `json:"createdAt"`
//	UpdatedAt time.Time `json:"updatedAt"`
//}

func (Play) IsNode() {}

type Play struct {
	Key       *datastore.Key `datastore:"__key__"`
	ID        int64          `json:"id" datastore:"-"`
	UserKey    *datastore.Key `json:"user_key" datastore:"user_key"`
	AudioKey   *datastore.Key `json:"audio_key" datastore:"audio_key"`
	CreatedAt   time.Time      `json:"created_at" datastore:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" datastore:"updated_at"`
}

func (a Play) GetKey() *datastore.Key {
	return datastore.IncompleteKey(AudioUserKind, nil)
}

func (a Play) SetID(key *datastore.Key) {
	a.ID = key.ID
}

func NewPlay(userID int64, audioID string) *Play {
	audioKey := GetAudioKey(audioID)
	userKey := GetUserKey(userID)
	au := &Play{
		UserKey:  userKey,
		AudioKey: audioKey,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return au
}
