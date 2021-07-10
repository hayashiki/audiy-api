package ds

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/hayashiki/audiy-api/domain/entity"
)

// AudioRepository operates Audio entity
type playRepository struct {
	client *datastore.Client
}

// Find finds audio given id
func (repo *playRepository) Find(ctx context.Context, userID int64, audioID string) (*entity.Play, error) {
	panic("implement me")
}

func NewPlayRepository(client *datastore.Client) entity.PlayRepository {
	return &playRepository{client: client}
}

func (repo *playRepository) Exists(ctx context.Context, userID int64, audioID string) (bool, error) {
	userKey := entity.GetUserKey(userID)
	audioKey := entity.GetAudioKey(audioID)
	q := datastore.NewQuery(entity.AudioUserKind).Filter("user_key=", userKey).Filter("audio_key=", audioKey).KeysOnly().Limit(1)
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
	_, err := repo.client.Put(ctx, datastore.IncompleteKey(entity.AudioUserKind, nil), item)
	return err
}