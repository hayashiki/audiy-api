package entity

import (
	"context"
	"github.com/hayashiki/audiy-api/infrastructure/ds"
)

// AudioUserRepository interface
type AudioUserRepository interface {
	Exists(context.Context, string) bool
	Find(context.Context, string) (*AudioUser, error)
	Save(context.Context, *AudioUser) error
	Remove(ctx context.Context, item *AudioUser) error
}

// AudioRepository operates Audio entity
type audioUserRepository struct {
	*ds.DataStore
}

// Find finds audio given id
func (repo *audioUserRepository) Find(ctx context.Context, id int64) (*AudioUser, error) {
	item := &AudioUser{ID: id}
	return item, repo.Get(ctx, item)
}

// Save saves audios
func (repo *audioUserRepository) Save(ctx context.Context, item *AudioUser) error {
	// TODO: if exists
	_, err := repo.Put(ctx, item)
	return err
}
