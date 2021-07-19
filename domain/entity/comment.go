package entity

import (
	"time"

	"cloud.google.com/go/datastore"
)

const CommentKind = "Comment"

type Comment struct {
	Key       *datastore.Key `datastore:"__key__"`
	ID        int64          `json:"id" datastore:"-"`
	UserKey   *datastore.Key `json:"user_key" datastore:"user_key"`
	AudioKey  *datastore.Key `json:"audio_key" datastore:"audio_key"`
	Body      string         `json:"body" datastore:"body"`
	CreatedAt time.Time      `json:"createdAt" datastore:"created_at"`
	UpdatedAt time.Time      `json:"updatedAt" datastore:"updated_at"`
}

func (Comment) IsNode() {}

func NewComment(userID string, audioID string, body string) *Comment {
	audioKey := GetAudioKey(audioID)
	userKey := GetUserKey(userID)
	au := &Comment{
		UserKey:   userKey,
		AudioKey:  audioKey,
		Body:      body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return au
}
