package ds

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/hayashiki/audiy-api/domain/entity"
)

// playRepository operates play entity
type likeRepository struct {
	client *datastore.Client
}

// Find finds play given userID and audioID
func (repo *likeRepository) Find(ctx context.Context, userID int64, audioID string) (*entity.Like, error) {
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
