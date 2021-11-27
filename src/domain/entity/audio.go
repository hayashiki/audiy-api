package entity

import (
	"fmt"
	"time"

	"cloud.google.com/go/datastore"
)

// AudioKind is kind name for audio app
const AudioKind = "Audio"

// Audio is an object representing the radio schema
type Audio struct {
	Key    *datastore.Key `datastore:"__key__"`
	ID     string         `json:"id" datastore:"-"`
	Name   string         `json:"name" datastore:"name"`
	Length float64            `json:"length" datastore:"length"`
	//URL         string         `json:"url" datastore:"url"`
	Mimetype     string    `json:"mimetype" datastore:"mimetype"`
	PublishedAt  time.Time `json:"published_at" datastore:"published_at"`
	CreatedAt    time.Time `json:"created_at" datastore:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" datastore:"updated_at"`
	LikeCount    int       `json:"like_count" datastore:"like_count"`
	PlayCount    int       `json:"play_count" datastore:"play_count"`
	CommentCount int       `json:"comment_count" datastore:"comment_count"`
}

func (a *Audio) GetKey() *datastore.Key {
	if a.ID == "" {
		return nil
	}
	return datastore.NameKey(AudioKind, a.ID, nil)
}

func (a *Audio) GetID() string {
	return a.ID
}

func (a *Audio) GetName() string {
	return a.Name
}

func (a *Audio) GetLength() float64 {
	return a.Length
}

func (a *Audio) GetMimetype() string {
	return a.Mimetype
}

func (a *Audio) GetCreatedAt() time.Time {
	return a.CreatedAt
}

func (a *Audio) GetUpdatedAt() time.Time {
	return a.UpdatedAt
}

// TODO
func (a *Audio) IsPublished() bool {
	return true
}

func GetAudioKey(id string) *datastore.Key {
	entity := Audio{ID: id}
	return entity.GetKey()
}

func NewAudio(id, name string, length float64, url, mimetype string, created time.Time) *Audio {
	return &Audio{
		ID:     id,
		Name:   name,
		Length: length,
		//URL:         url,
		Mimetype:    mimetype,
		PublishedAt: created,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (a *Audio) SetID(key *datastore.Key) {
	a.ID = key.Name
}

func (Audio) IsNode() {}

func (a *Audio) String() string {
	return fmt.Sprintf("ID: %v, Name: %s", a.ID, a.Name)
}

// TODO: Validate
func (a *Audio) Validate() error {
	return nil
}

// TODO Validation helper to check mimetype
func validMimetype(value interface{}) error {
	return nil
}
