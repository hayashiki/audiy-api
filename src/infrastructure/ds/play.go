package ds

import (
	"context"
	"fmt"

	entity2 "github.com/hayashiki/audiy-api/src/domain/entity"

	"cloud.google.com/go/datastore"
)

// playRepository operates play entity
type playRepository struct {
	client *datastore.Client
}

// Find finds play given userID and audioID
func (repo *playRepository) Find(ctx context.Context, userID string, audioID string) (*entity2.Play, error) {
	panic("implement me")
}

func NewPlayRepository(client *datastore.Client) entity2.PlayRepository {
	return &playRepository{client: client}
}

func (repo *playRepository) Exists(ctx context.Context, userID string, audioID string) (bool, error) {
	userKey := entity2.GetUserKey(userID)
	audioKey := entity2.GetAudioKey(audioID)
	q := datastore.NewQuery(entity2.PlayKind).Filter("user_key=", userKey).Filter("audio_key=", audioKey).KeysOnly().Limit(1)
	var dst []*entity2.Play

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
func (repo *playRepository) Save(ctx context.Context, item *entity2.Play) error {
	// TODO: if exists
	key, err := repo.client.Put(ctx, datastore.IncompleteKey(entity2.PlayKind, nil), item)
	item.ID = key.ID
	return err
}
