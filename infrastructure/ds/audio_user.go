package ds

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/hayashiki/audiy-api/domain/entity"
)

// AudioRepository operates Audio entity
type audioUserRepository struct {
	client *datastore.Client
}

// Find finds audio given id
func (repo *audioUserRepository) Find(ctx context.Context, userID int64, audioID string) (*entity.AudioUser, error) {
	panic("implement me")
}

func NewAudioUserRepository(client *datastore.Client) entity.AudioUserRepository {
	return &audioUserRepository{client: client}
}

func (repo *audioUserRepository) Exists(ctx context.Context, userID int64, audioID string) (bool, error) {
	q := datastore.NewQuery(entity.AudioUserKind).Filter("user_id=", userID).Filter("audio_id=", audioID).KeysOnly().Limit(1)
	var dst []*entity.UserAudio

	keys, err := repo.client.GetAll(ctx, q, dst)
	if err != nil {
		return false, fmt.Errorf("not found user %w", err)
	}
	if len(keys) == 1 {
		return true, nil
	}
	return false, nil
}

// Save saves audios
func (repo *audioUserRepository) Save(ctx context.Context, item *entity.AudioUser) error {
	// TODO: if exists
	_, err := repo.client.Put(ctx,  datastore.IncompleteKey(entity.AudioUserKind, nil),item)
	return err
}
