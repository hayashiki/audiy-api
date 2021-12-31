package model

import (
	"fmt"
	"time"
)

const FeedKind = "Feed"

type FeedID string

type Feed struct {
	//ID          string          `json:"id"`
	AudioID     string
	UserID      string
	Played      bool           `json:"played"`
	Liked       bool           `json:"liked"`
	Stared      bool           `json:"stared"`
	StartTime   float64       `json:"start_time"`
	PublishedAt time.Time      `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (Feed) IsNode() {}

func NewFeed(audioID string, userID string, publishedAt time.Time) *Feed {
	return &Feed{
		//ID: ID(audioID, userID),
		PublishedAt: publishedAt,
		AudioID:     audioID,
		UserID:      userID,
		Played:      false,
		Liked:       false,
		Stared:      false,
		StartTime:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (f *Feed) ID() FeedID {
	return FeedID(fmt.Sprintf("%s-%s", f.AudioID, f.UserID))
}

func NewFeedID(audioID string, userID string) FeedID {
	return (&Feed{AudioID: audioID, UserID: userID}).ID()
}

func (f *Feed) Like() {
	f.Liked = true
}

func (f *Feed) UnLike() {
	f.Liked = false
}

func (f *Feed) Play() {
	f.Played = true
}

func (f *Feed) UnPlay() {
	f.Played = false
}

func (f *Feed) Star() {
	f.Stared = true
}

func (f *Feed) UnStar() {
	f.Stared = false
}

