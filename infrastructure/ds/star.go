package ds

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/hayashiki/audiy-api/domain/entity"
)

// playRepository operates play entity
type starRepository struct {
	client *datastore.Client
}

// Find finds star by id
func (repo *starRepository) Find(ctx context.Context, id int64) (*entity.Star, error) {
	var dst *entity.Star
	key := datastore.IDKey(entity.StarKind, id, nil)
	err := repo.client.Get(ctx, key, &dst)

	dst.ID = key.ID
	return dst, err
}

// FindByRel finds star given userID and audioID
func (repo *starRepository) FindByRel(ctx context.Context, userID string, audioID string) (*entity.Star, error) {
	userKey := entity.GetUserKey(userID)
	audioKey := entity.GetAudioKey(audioID)
	q := datastore.NewQuery(entity.StarKind).Filter("user_key=", userKey).Filter("audio_key=", audioKey).KeysOnly().Limit(1)
	var dst []*entity.Star

	keys, err := repo.client.GetAll(ctx, q, dst)
	if err != nil {
		return nil, fmt.Errorf("not found user %w", err)
	}
	dst[0].ID = keys[0].ID
	return dst[0], err
}

func NewStarRepository(client *datastore.Client) entity.StarRepository {
	return &starRepository{client: client}
}

func (repo *starRepository) Exists(ctx context.Context, userID string, audioID string) (bool, error) {
	userKey := entity.GetUserKey(userID)
	audioKey := entity.GetAudioKey(audioID)
	q := datastore.NewQuery(entity.StarKind).Filter("user_key=", userKey).Filter("audio_key=", audioKey).KeysOnly().Limit(1)
	var dst []*entity.Star

	keys, err := repo.client.GetAll(ctx, q, dst)
	if err != nil {
		return false, fmt.Errorf("not found user %w", err)
	}
	if len(keys) == 1 {
		return true, nil
	}
	return false, nil
}

// Save saves star
func (repo *starRepository) Save(ctx context.Context, item *entity.Star) error {
	// TODO: if exists
	key, err := repo.client.Put(ctx, datastore.IncompleteKey(entity.StarKind, nil), item)
	item.ID = key.ID
	return err
}

// Delete star
func (repo *starRepository) Delete(ctx context.Context, id int64) error {
	err := repo.client.Delete(ctx, datastore.IDKey(entity.StarKind, id, nil))
	return err
}
