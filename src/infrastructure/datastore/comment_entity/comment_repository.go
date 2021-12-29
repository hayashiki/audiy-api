package comment_entity

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

func NewCommentRepository(client datastore.Client) repository.CommentRepository {
	return &repo{
		client: client,
	}
}

func (r *repo) GetAllByAudio(
	ctx context.Context,
	audioID string,
	cursor string,
	limit int,
	orderBy string) ([]*model.Comment, string, bool, error) {

	filters := map[string]interface{}{
		"AudioID=": audioID,
	}

	keys, nextCursor, hasMore, err := r.client.RunQuery(
		ctx, kind, filters, cursor, limit, orderBy)
	if err != nil {
		return nil, nextCursor, hasMore, errors.WithStack(err)
	}
	entities := make([]*entity, 0, len(keys))
	for _, id := range keys {
		entities = append(entities, onlyID(id.ID()))
	}

	if err := r.client.GetMulti(ctx, entities); err != nil {
		return nil, nextCursor, hasMore, errors.WithStack(err)
	}
	comments := make([]*model.Comment, len(entities))
	for i, e := range entities {
		comments[i] = e.toDomain()
	}
	return comments, nextCursor, hasMore, err
}

func (r *repo) Get(ctx context.Context, id int64) (*model.Comment, error) {
	entity := onlyID(id)

	if err := r.client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.toDomain(), nil
}

func (r *repo) Put(ctx context.Context, item *model.Comment) error {
	if err := r.client.Put(ctx, toEntity(item)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repo) PutTx(tx *boom.Transaction, item *model.Comment) error {
	if err := r.client.PutTx(tx, toEntity(item)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// TODO: idに型をつけよう。。
func (r *repo) DeleteTx(tx *boom.Transaction, id int64) error {
	if err := r.client.DeleteTx(tx, onlyID(id)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	if err := r.client.Delete(ctx, onlyID(id)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

