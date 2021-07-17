package ds

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/hayashiki/audiy-api/domain/entity"
)

// playRepository operates play entity
type likeRepository struct {
	client *datastore.Client
}

func (repo *likeRepository) FindByRel(ctx context.Context, userID int64, audioID string) (*entity.Like, error) {
	userKey := entity.GetUserKey(userID)
	audioKey := entity.GetAudioKey(audioID)
	q := datastore.NewQuery(entity.PlayKind).Filter("user_key=", userKey).Filter("audio_key=", audioKey).Limit(1)
	var dst []*entity.Like

	keys, err := repo.client.GetAll(ctx, q, dst)
	if err != nil {
		return nil, fmt.Errorf("not found like %w", err)
	}

	dst[0].ID = keys[0].ID
	return dst[0], nil
}

func (repo *likeRepository) Delete(ctx context.Context, id int64) error {
	key := datastore.IDKey(entity.LikeKind, id, nil)
	return repo.client.Delete(ctx, key)
}

// Find finds play given userID and audioID
func (repo *likeRepository) Find(ctx context.Context, id int64) (*entity.Like, error) {
	var dst *entity.Like
	key := datastore.IDKey(entity.LikeKind, id, nil)

	err := repo.client.Get(ctx, key, &dst)
	if err != nil {
		return nil, fmt.Errorf("not found like %w", err)
	}
	dst.ID = key.ID
	return dst, nil
}

func NewLikeRepository(client *datastore.Client) entity.LikeRepository {
	return &likeRepository{client: client}
}

func (repo *likeRepository) Exists(ctx context.Context, userID int64, audioID string) (bool, error) {
	userKey := entity.GetUserKey(userID)
	audioKey := entity.GetAudioKey(audioID)
	q := datastore.NewQuery(entity.LikeKind).Filter("user_key=", userKey).Filter("audio_key=", audioKey).KeysOnly().Limit(1)
	var dst []*entity.Like

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
func (repo *likeRepository) Save(ctx context.Context, item *entity.Like) error {
	// TODO: if exists
	_, err := repo.client.Put(ctx, datastore.IncompleteKey(entity.PlayKind, nil), item)
	return err
}
