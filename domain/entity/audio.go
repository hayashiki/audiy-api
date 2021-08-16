package entity

import (
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
	Length int            `json:"length" datastore:"length"`
	//URL         string         `json:"url" datastore:"url"`
	Mimetype    string    `json:"mimetype" datastore:"mimetype"`
	PublishedAt time.Time `json:"published_at" datastore:"published_at"`
	CreatedAt   time.Time `json:"created_at" datastore:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" datastore:"updated_at"`
	PlayedUsers []string  `json:"played_users" datastore:"played_users"`
}

func (r *Audio) GetKey() *datastore.Key {
	if r.ID == "" {
		return nil
	}
	return datastore.NameKey(AudioKind, r.ID, nil)
}

func GetAudioKey(id string) *datastore.Key {
	entity := Audio{ID: id}
	return entity.GetKey()
}

func NewAudio(id, name string, length int, url, mimetype string, created time.Time) *Audio {
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

func (r *Audio) SetID(key *datastore.Key) {
	r.ID = key.Name
}

func (Audio) IsNode() {}
