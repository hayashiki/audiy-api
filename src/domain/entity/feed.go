package entity

import (
	"time"

	"cloud.google.com/go/datastore"
)

const FeedKind = "Feed"

type Feed struct {
	Key         *datastore.Key `datastore:"__key__"`
	ID          int64          `json:"id" datastore:"-"`
	PublishedAt time.Time      `json:"published_at" datastore:"published_at"`
	AudioKey    *datastore.Key `json:"audio_key" datastore:"audio_key"`
	Played      bool           `json:"played" datastore:"played"`
	Liked       bool           `json:"liked" datastore:"liked"`
	Stared      bool           `json:"stared" datastore:"stared"`
	StartTime   *float64       `json:"start_time" datastore:"start_time"`
	CreatedAt   time.Time      `json:"created_at" datastore:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" datastore:"updated_at"`
}

func (Feed) IsNode() {}

func NewFeed(audioID string, publishedAt time.Time) *Feed {
	audioKey := GetAudioKey(audioID)

	return &Feed{
		PublishedAt: publishedAt,
		AudioKey:    audioKey,
		Played:      false,
		Liked:       false,
		Stared:      false,
		StartTime:   nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
