package ds

import (
	"context"
	"fmt"

	entity2 "github.com/hayashiki/audiy-api/src/domain/entity"

	"cloud.google.com/go/datastore"
)

// playRepository operates play entity
type likeRepository struct {
	client *datastore.Client
}

func (repo *likeRepository) FindByRel(ctx context.Context, userID string, audioID string) (*entity2.Like, error) {
	userKey := entity2.GetUserKey(userID)
	audioKey := entity2.GetAudioKey(audioID)
	q := datastore.NewQuery(entity2.PlayKind).Filter("user_key=", userKey).Filter("audio_key=", audioKey).Limit(1)
	var dst []*entity2.Like

	keys, err := repo.client.GetAll(ctx, q, dst)
	if err != nil {
		return nil, fmt.Errorf("not found like %w", err)
	}

	dst[0].ID = keys[0].ID
	return dst[0], nil
}

func (repo *likeRepository) Delete(ctx context.Context, id int64) error {
	key := datastore.IDKey(entity2.LikeKind, id, nil)
	return repo.client.Delete(ctx, key)
}

// Find finds play given userID and audioID
func (repo *likeRepository) Find(ctx context.Context, id int64) (*entity2.Like, error) {
	var dst *entity2.Like
	key := datastore.IDKey(entity2.LikeKind, id, nil)

	err := repo.client.Get(ctx, key, &dst)
	if err != nil {
		return nil, fmt.Errorf("not found like %w", err)
	}
	dst.ID = key.ID
	return dst, nil
}

func NewLikeRepository(client *datastore.Client) entity2.LikeRepository {
	return &likeRepository{client: client}
}

func (repo *likeRepository) Exists(ctx context.Context, userID string, audioID string) (bool, error) {
	userKey := entity2.GetUserKey(userID)
	audioKey := entity2.GetAudioKey(audioID)
	q := datastore.NewQuery(entity2.LikeKind).Filter("user_key=", userKey).Filter("audio_key=", audioKey).KeysOnly().Limit(1)
	var dst []*entity2.Like

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
func (repo *likeRepository) Save(ctx context.Context, item *entity2.Like) error {
	// TODO: if exists
	key, err := repo.client.Put(ctx, datastore.IncompleteKey(entity2.LikeKind, nil), item)
	item.ID = key.ID
	return err
}
