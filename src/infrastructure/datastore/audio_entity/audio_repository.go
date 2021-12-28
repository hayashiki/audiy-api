package audio_entity

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

func NewAudioRepository(client datastore.Client) repository.AudioRepository {
	return &repo{
		client: client,
	}
}

func (r *repo) GetAll(
	ctx context.Context,
	cursor string,
	limit int,
	orderBy string) ([]*model.Audio, string, bool, error) {

	filters := map[string]interface{}{
		//"AudioID=": audioID,
	}
	//
	keys, nextCursor, hasMore, err := r.client.RunQuery(
		ctx, kind, filters, cursor, limit, orderBy)
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
	audios := make([]*model.Audio, len(entities))
	for i, e := range entities {
		audios[i] = e.toDomain()
	}
	return audios, nextCursor, hasMore, err
}

func (r *repo) GetMulti(ctx context.Context, ids []string) ([]*model.Audio, error) {
	entities := make([]*entity, 0, len(ids))
	for _, id := range ids {
		entities = append(entities, onlyID(id))
	}

	if err := r.client.GetMulti(ctx, entities); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*model.Audio, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repo) Get(ctx context.Context, id string) (*model.Audio, error) {
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

func (r *repo) Put(ctx context.Context, item *model.Audio) error {
	if err := r.client.Put(ctx, toEntity(item)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repo) PutTx(tx *boom.Transaction, item *model.Audio) error {
	if err := r.client.PutTx(tx, toEntity(item)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// TODO: idに型をつけよう。。
func (r *repo) DeleteTx(tx *boom.Transaction, id string) error {
	if err := r.client.DeleteTx(tx, onlyID(id)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}


