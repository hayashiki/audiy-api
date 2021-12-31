package fcm_entity

import (
	clouddatastore "cloud.google.com/go/datastore"
	"context"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/domain/repository"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
	"github.com/pkg/errors"
)

type repo struct {
	client datastore.DSClient
}

func NewRepository(client datastore.DSClient) repository.FCMRepository {
	return &repo{
		client: client,
	}
}

func (r *repo) GetAll(
	ctx context.Context,
	cursor string,
	limit int,
	orderBy string) ([]*model.Fcm, string, bool, error) {
	q := &model.Query{
		Kind: kind,
		Limit:     limit,
		Cursor:    cursor,
		Filters:   nil,
		OrderBy:   orderBy,
		Namespace: "",
	}
	keys, nextCursor, hasMore, err := r.client.Run(ctx, q)
	if err != nil {
		return nil, nextCursor, hasMore, errors.WithStack(err)
	}
	entities := make([]*entity, 0, len(keys))
	if err := r.client.GetMulti(ctx, keys, entities); err != nil {
		return nil, nextCursor, hasMore, errors.WithStack(err)
	}
	fcms := make([]*model.Fcm, len(entities))
	for i, e := range entities {
		fcms[i] = e.toDomain()
	}
	return fcms, nextCursor, hasMore, err
}

func (r *repo) Get(ctx context.Context, id string) (*model.Fcm, error) {
	item := &entity{}
	key := clouddatastore.NameKey(kind, id, nil)

	if err := r.client.Get(ctx, key, item); err != nil {
		return nil, errors.WithStack(err)
	}

	return item.toDomain(), nil
}

func (r *repo) PutTx(tx *clouddatastore.Transaction, item *model.Fcm) error {
	key := clouddatastore.NameKey(kind, item.ID, nil)

	if err := r.client.PutTx(tx, key, toEntity(item)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repo) Delete(tx *clouddatastore.Transaction, id string) error {
	key := clouddatastore.NameKey(kind, id, nil)

	if err := r.client.DeleteTx(tx, key); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
