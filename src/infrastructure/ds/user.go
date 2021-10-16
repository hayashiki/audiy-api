package ds

import (
	"context"
	"errors"
	"fmt"

	entity2 "github.com/hayashiki/audiy-api/src/domain/entity"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

// AudioRepository operates Audio entity
type userRepository struct {
	client *datastore.Client
}

func NewUserRepository(client *datastore.Client) entity2.UserRepository {
	return &userRepository{client: client}
}

func (repo *userRepository) Exists(ctx context.Context, userKey *datastore.Key) (bool, error) {
	q := datastore.NewQuery(entity2.UserKind).Filter("__key__=", userKey).KeysOnly().Limit(1)
	var dst []*entity2.User

	keys, err := repo.client.GetAll(ctx, q, dst)
	if err != nil {
		return false, fmt.Errorf("not found user %w", err)
	}
	if len(keys) == 1 {
		return true, nil
	}
	return false, nil
}

// Get user
func (repo *userRepository) Get(ctx context.Context, userID string) (*entity2.User, error) {
	var user entity2.User
	err := repo.client.Get(ctx, datastore.NameKey(entity2.UserKind, userID, nil), &user)
	if err != nil {
		return nil, err
	}
	user.ID = user.Key.Name
	return &user, err
}

// Get users
func (repo *userRepository) GetAll(ctx context.Context) ([]*entity2.User, error) {
	query := datastore.NewQuery(entity2.UserKind)
	it := repo.client.Run(ctx, query)
	entities := make([]*entity2.User, 0)
	for {
		entity := &entity2.User{}
		_, err := it.Next(entity)
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return entities, err
		}
		entity.ID = entity.Key.Name
		entities = append(entities, entity)
	}

	return entities, nil
}

// Save saves user
func (repo *userRepository) Save(ctx context.Context, user *entity2.User) error {
	// TODO: if exists
	newKey := datastore.NameKey(entity2.UserKind, user.ID, nil)
	key, err := repo.client.Put(ctx, newKey, user)
	user.Key = key
	return err
}

// Delete deletes user
func (repo *userRepository) Delete(ctx context.Context, userKey *datastore.Key) error {
	err := repo.client.Delete(ctx, userKey)
	return err
}
