package user_entity

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/domain/repository"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
	"github.com/pkg/errors"
	"go.mercari.io/datastore/boom"
)

type repo struct {
	client datastore.Client
}

func NewUserRepository(client datastore.Client) repository.UserRepository {
	return &repo{
		client: client,
	}
}

func (r *repo) GetAll(ctx context.Context) ([]*model.User, error) {
	var entities []*entity
	if err := r.client.GetAll(ctx, kind, &entities); err != nil {
		return nil, err
	}
	users := make([]*model.User, len(entities))
	for i, e := range entities {
		users[i] = e.toDomain()
	}
	return users, nil
}

func (r *repo) Get(ctx context.Context, id string) (*model.User, error) {
	entity := onlyID(id)

	if err := r.client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.toDomain(), nil
}

func (r *repo) Exists(ctx context.Context, id string) (bool, error) {
	exists, err := r.client.Exists(ctx, onlyID(id))
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exists, nil
}

func (r *repo) Put(ctx context.Context, item *model.User) error {
	if err := r.client.Put(ctx, toEntity(item)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repo) PutTx(tx *boom.Transaction, item *model.User) error {
	if err := r.client.PutTx(tx, toEntity(item)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repo) DeleteTx(tx *boom.Transaction, id string) error {
	if err := r.client.DeleteTx(tx, onlyID(id)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	if err := r.client.Delete(ctx, onlyID(id)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

