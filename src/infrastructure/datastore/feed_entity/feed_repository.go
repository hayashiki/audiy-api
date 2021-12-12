package feed_entity

import (
	clouddatastore "cloud.google.com/go/datastore"
	"context"
	"github.com/hayashiki/audiy-api/src/domain/repository"
	"log"

	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
	"github.com/pkg/errors"
)

type repo struct {
	client datastore.DSClient
}

func NewFeedRepository(client datastore.DSClient) repository.FeedRepository {
	return &repo{
		client: client,
	}
}

func (r *repo) GetAll(
	ctx context.Context,
	userID string,
	filters map[string]interface{},
	cursor string,
	limit int,
	orderBy string) ([]*model.Feed, string, bool, error) {

	var keys []*clouddatastore.Key
	var parentKey *clouddatastore.Key
	parentKey = clouddatastore.NameKey(parentKind, userID, parentKey)
	//filters := map[string]interface{}{
	//	//"AudioID=": audioID,
	//}
	keys, nextCursor, hasMore, err := r.client.RunQuery(
		ctx, kind, parentKey, filters, cursor, limit, orderBy)
	if err != nil {
		return nil, nextCursor, hasMore, errors.WithStack(err)
	}
	entities := make([]*entity, len(keys))
	if err := r.client.GetMulti(ctx, keys, entities); err != nil {
		return nil, nextCursor, hasMore, errors.WithStack(err)
	}
	audios := make([]*model.Feed, len(entities))
	for i, e := range entities {
		audios[i] = e.toDomain()
	}
	return audios, nextCursor, hasMore, err
}

func (r *repo) GetMulti(ctx context.Context, userID string, ids []int64) ([]*model.Feed, error) {

	keys := make([]*clouddatastore.Key, 0, len(ids))
	var parentKey *clouddatastore.Key
	parentKey = clouddatastore.NameKey("User", userID, nil)

	entities := make([]*entity, 0, len(ids))
	for i, id := range ids {
		entities[i] = onlyID(id)
		keys[i] = clouddatastore.IDKey(kind, id, parentKey)
	}

	if err := r.client.GetMulti(ctx, keys, entities); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*model.Feed, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repo) Get(ctx context.Context, userID string, id int64) (*model.Feed, error) {
	entity := onlyID(id)
	parentKey := clouddatastore.NameKey(parentKind, userID, nil)
	key := clouddatastore.IDKey(kind, id, parentKey)
	if err := r.client.Get(ctx, key, entity); err != nil {
		log.Println(entity, err)
		return nil, errors.WithStack(err)
	}

	return entity.toDomain(), nil
}

// TODO: ancestor
func (r *repo) GetByAudio(ctx context.Context, userID string, audioID string) (*model.Feed, error) {
	filters := map[string]interface{}{
		"AudioID=": audioID,
	}
	parentKey := clouddatastore.NameKey(parentKind, userID, nil)
	keys, _, _, err := r.client.RunQuery(ctx, kind, parentKey, filters, "", 1, "")
	if err != nil {
		return nil, err
	}
	var item *entity

	if err := r.client.Get(ctx, keys[0], &item); err != nil {
		return nil, err
	}

	return item.toDomain(), nil
}

func (r *repo) Exists(ctx context.Context, userID string, id int64) (bool, error) {
	parentKey := clouddatastore.NameKey(parentKind, userID, nil)
	key := clouddatastore.IDKey(kind, id, parentKey)
	exists, err := r.client.Exists(ctx, key, onlyID(id))
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exists, nil

}

func (r *repo) Put(ctx context.Context, userID string, item *model.Feed) error {
	parentKey := clouddatastore.NameKey(parentKind, userID, nil)
	key := clouddatastore.IDKey(kind, item.ID, parentKey)
	if err := r.client.Put(ctx, key, toEntity(item)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repo) PutMulti(ctx context.Context, feeds []*model.Feed) error {
	keys := make([]*clouddatastore.Key, len(feeds))

	for i, item := range feeds {
		parentKey := clouddatastore.NameKey(kind, item.UserID, nil)
		keys[i] = clouddatastore.IDKey(kind, item.ID, parentKey)
	}

	err := r.client.PutMulti(ctx, keys, feeds)
	//feed.Key = key

	return err
}

func (r *repo) PutTx(tx *clouddatastore.Transaction, userID string, item *model.Feed) error {
	parentKey := clouddatastore.NameKey(parentKind, userID, nil)
	key := clouddatastore.IDKey(kind, item.ID, parentKey)

	if err := r.client.PutTx(tx, key, toEntity(item)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
//
//// TODO: idに型をつけよう。。
func (r *repo) DeleteTx(tx *clouddatastore.Transaction, userID string, id int64) error {
	pKey := clouddatastore.NameKey(parentKind, userID, nil)
	key := clouddatastore.IDKey(kind, id, pKey)

	if err := r.client.DeleteTx(tx, key); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
