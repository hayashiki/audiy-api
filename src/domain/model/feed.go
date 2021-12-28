package model

import (
	"time"

	"cloud.google.com/go/datastore"
)

const FeedKind = "Feed"

type Feed struct {
	Key         *datastore.Key `datastore:"__key__"`
	ID          int64          `json:"id" datastore:"-"`
	AudioID     string
	UserID      string
	Played      bool           `json:"played"`
	Liked       bool           `json:"liked"`
	Stared      bool           `json:"stared"`
	StartTime   *float64       `json:"start_time"`
	PublishedAt time.Time      `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (Feed) IsNode() {}

func NewFeed(audioID string, userID string, publishedAt time.Time) *Feed {
	return &Feed{
		ID: NewID("Feed"),
		PublishedAt: publishedAt,
		AudioID:     audioID,
		UserID:      userID,
		Played:      false,
		Liked:       false,
		Stared:      false,
		StartTime:   nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
