package fcm_entity

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

func NewRepository(client datastore.Client) repository.FCMRepository {
	return &repo{
		client: client,
	}
}

func (r *repo) GetAll(
	ctx context.Context,
	cursor string,
	limit int,
	orderBy string) ([]*model.Fcm, string, bool, error) {
	keys, nextCursor, hasMore, err := r.client.RunQuery(
		ctx, kind, nil, cursor, limit, orderBy)
	if err != nil {
		return nil, nextCursor, hasMore, errors.WithStack(err)
	}
	entities := make([]*entity, 0, len(keys))
	for _, id := range keys {
		entities = append(entities, onlyID(id.Name()))
	}

	if err := r.client.GetMulti(ctx, entities); err != nil {
		return nil, nextCursor, hasMore, errors.WithStack(err)
	}
	fcms := make([]*model.Fcm, len(entities))
	for i, e := range entities {
		fcms[i] = e.toDomain()
	}
	return fcms, nextCursor, hasMore, err
}

func (r *repo) Get(ctx context.Context, id string) (*model.Fcm, error) {
	entity := onlyID(id)

	if err := r.client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.toDomain(), nil
}

func (r *repo) Put(tx *boom.Transaction, item *model.Fcm) error {
	if err := r.client.PutTx(tx, toEntity(item)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// TODO: idに型をつけよう。。
func (r *repo) Delete(tx *boom.Transaction, id string) error {
	if err := r.client.DeleteTx(tx, onlyID(id)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
