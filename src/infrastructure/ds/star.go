package ds

import (
	"context"
	"fmt"

	entity2 "github.com/hayashiki/audiy-api/src/domain/entity"

	"cloud.google.com/go/datastore"
)

// playRepository operates play entity
type starRepository struct {
	client *datastore.Client
}

// Find finds star by id
func (repo *starRepository) Find(ctx context.Context, id int64) (*entity2.Star, error) {
	var dst *entity2.Star
	key := datastore.IDKey(entity2.StarKind, id, nil)
	err := repo.client.Get(ctx, key, &dst)

	dst.ID = key.ID
	return dst, err
}

// FindByRel finds star given userID and audioID
func (repo *starRepository) FindByRel(ctx context.Context, userID string, audioID string) (*entity2.Star, error) {
	userKey := entity2.GetUserKey(userID)
	audioKey := entity2.GetAudioKey(audioID)
	q := datastore.NewQuery(entity2.StarKind).Filter("user_key=", userKey).Filter("audio_key=", audioKey).KeysOnly().Limit(1)
	var dst []*entity2.Star

	keys, err := repo.client.GetAll(ctx, q, dst)
	if err != nil {
		return nil, fmt.Errorf("not found user %w", err)
	}
	dst[0].ID = keys[0].ID
	return dst[0], err
}

func NewStarRepository(client *datastore.Client) entity2.StarRepository {
	return &starRepository{client: client}
}

func (repo *starRepository) Exists(ctx context.Context, userID string, audioID string) (bool, error) {
	userKey := entity2.GetUserKey(userID)
	audioKey := entity2.GetAudioKey(audioID)
	q := datastore.NewQuery(entity2.StarKind).Filter("user_key=", userKey).Filter("audio_key=", audioKey).KeysOnly().Limit(1)
	var dst []*entity2.Star

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
func (repo *starRepository) Save(ctx context.Context, item *entity2.Star) error {
	// TODO: if exists
	key, err := repo.client.Put(ctx, datastore.IncompleteKey(entity2.StarKind, nil), item)
	item.ID = key.ID
	return err
}

// Delete star
func (repo *starRepository) Delete(ctx context.Context, id int64) error {
	err := repo.client.Delete(ctx, datastore.IDKey(entity2.StarKind, id, nil))
	return err
}
