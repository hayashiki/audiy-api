package ds

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/hayashiki/audiy-api/domain/entity"
)

// playRepository operates play entity
type playRepository struct {
	client *datastore.Client
}

// Find finds play given userID and audioID
func (repo *playRepository) Find(ctx context.Context, userID string, audioID string) (*entity.Play, error) {
	panic("implement me")
}

func NewPlayRepository(client *datastore.Client) entity.PlayRepository {
	return &playRepository{client: client}
}

func (repo *playRepository) Exists(ctx context.Context, userID string, audioID string) (bool, error) {
	userKey := entity.GetUserKey(userID)
	audioKey := entity.GetAudioKey(audioID)
	q := datastore.NewQuery(entity.PlayKind).Filter("user_key=", userKey).Filter("audio_key=", audioKey).KeysOnly().Limit(1)
	var dst []*entity.Play

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
func (repo *playRepository) Save(ctx context.Context, item *entity.Play) error {
	// TODO: if exists
	key, err := repo.client.Put(ctx, datastore.IncompleteKey(entity.PlayKind, nil), item)
	item.ID = key.ID
	return err
}
