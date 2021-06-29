package entity

import (
	"context"
	"github.com/hayashiki/audiy-api/infrastructure/ds"
)

// AudioRepository interface
type AudioRepository interface {
	Exists(context.Context, string) bool
	Find(context.Context, string) (*Audio, error)
	FindAll(ctx context.Context, cursor string, limit int, sort ...string) ([]*Audio, string, error)
	Save(context.Context, *Audio) error
	Remove(ctx context.Context, item *Audio) error
}

// AudioRepository operates Audio entity
type audioRepository struct {
	*ds.DataStore
}

// NewAudioRepository returns the AudioRepository
func NewAudioRepository(store *ds.DataStore) AudioRepository {
	return &audioRepository{store}
}

// Exists exists item
func (repo *audioRepository) Exists(ctx context.Context, id string) bool {
	_, err := repo.Find(ctx, id)
	return err == nil
}

// FindAll finds all radios
func (repo *audioRepository) FindAll(ctx context.Context, cursor string, limit int, sort ...string) ([]*Audio, string, error) {
	q := ds.NewCursorQuery(AudioKind, nil, limit, cursor, sort...)
	rr, nextCursor, err := repo.GetAll(ctx, q, func() interface{} {
		var radio Audio
		return &radio
	})
	radios := make([]*Audio, len(rr))
	for i, r := range rr {
		ar := r.(*Audio)
		ar.SetID(ar.Key)
		radios[i] = ar
	}
	return radios, nextCursor, err
}

// Find finds audio given id
func (repo *audioRepository) Find(ctx context.Context, id string) (*Audio, error) {
	item := &Audio{ID: id}
	return item, repo.Get(ctx, item)
}

// Save saves audios
func (repo *audioRepository) Save(ctx context.Context, item *Audio) error {
	_, err := repo.Put(ctx, item)
	return err
}

// Delete saves audios
func (repo *audioRepository) Remove(ctx context.Context, item *Audio) error {
	_, err := repo.Delete(ctx, item)
	return err
}

