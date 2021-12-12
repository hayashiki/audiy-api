package model

import (
	"time"
)

type Comment struct {
	ID        int64          `json:"id" datastore:"-"`
	UserID    string
	AudioID    string
	Body      string         `json:"body"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

func (Comment) IsNode() {}

func NewComment(userID string, audioID string, body string) *Comment {
	comment := &Comment{
		ID: NewID("Comment"),
		AudioID:   audioID,
		UserID:    userID,
		Body:      body,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return comment
}
