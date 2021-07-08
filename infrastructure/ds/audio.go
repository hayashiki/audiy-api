package ds

import (
	"context"
	"github.com/hayashiki/audiy-api/domain/entity"
)

// AudioRepository operates Audio entity
type audioRepository struct {
	*DataStore
}

// NewAudioRepository returns the AudioRepository
func NewAudioRepository(store *DataStore) entity.AudioRepository {
	return &audioRepository{store}
}

// Exists exists item
func (repo *audioRepository) Exists(ctx context.Context, id string) bool {
	_, err := repo.Find(ctx, id)
	return err == nil
}

// FindAll finds all radios
func (repo *audioRepository) FindAll(ctx context.Context, cursor string, limit int, sort ...string) ([]*entity.Audio, string, error) {
	q := NewCursorQuery(entity.AudioKind, nil, limit, cursor, sort...)
	rr, nextCursor, err := repo.GetAll(ctx, q, func() interface{} {
		var radio entity.Audio
		return &radio
	})
	radios := make([]*entity.Audio, len(rr))
	for i, r := range rr {
		ar := r.(*entity.Audio)
		ar.SetID(ar.Key)
		radios[i] = ar
	}
	return radios, nextCursor, err
}

// Find finds audio given id
func (repo *audioRepository) Find(ctx context.Context, id string) (*entity.Audio, error) {
	item := &entity.Audio{ID: id}
	return item, repo.Get(ctx, item)
}

// Save saves audios
func (repo *audioRepository) Save(ctx context.Context, item *entity.Audio) error {
	_, err := repo.Put(ctx, item)
	return err
}

// Delete saves audios
func (repo *audioRepository) Remove(ctx context.Context, item *entity.Audio) error {
	_, err := repo.Delete(ctx, item)
	return err
}


