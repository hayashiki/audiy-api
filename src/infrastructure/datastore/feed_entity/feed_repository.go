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
	q := &model.Query{
		Kind: kind,
		Parent: parentKey,
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
	log.Println("keys", len(keys), keys[0])
	entities := make([]*entity, len(keys))
	if err := r.client.GetMulti(ctx, keys, entities); err != nil {
		log.Println("keys", len(keys), keys[0], err)
		return nil, nextCursor, hasMore, errors.WithStack(err)
	}
	feeds := make([]*model.Feed, len(entities))
	for i, e := range entities {
		feeds[i] = e.toDomain()
	}
	return feeds, nextCursor, hasMore, err
}

func (r *repo) GetMulti(ctx context.Context, userID string, ids []string) ([]*model.Feed, error) {

	keys := make([]*clouddatastore.Key, 0, len(ids))
	var parentKey *clouddatastore.Key
	parentKey = clouddatastore.NameKey("User", userID, nil)

	entities := make([]*entity, 0, len(ids))
	for i, id := range ids {
		entities[i] = onlyID(id)
		keys[i] = clouddatastore.NameKey(kind, id, parentKey)
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

func (r *repo) Get(ctx context.Context, id string, userID string) (*model.Feed, error) {
	en := &entity{}
	parentKey := clouddatastore.NameKey(parentKind, userID, nil)
	key := clouddatastore.NameKey(kind, id, parentKey)
	if err := r.client.Get(ctx, key, en); err != nil {
		return nil, errors.WithStack(err)
	}
	return en.toDomain(), nil
}

// TODO: ancestor
func (r *repo) GetByAudio(ctx context.Context, userID string, audioID string) (*model.Feed, error) {
	filters := map[string]interface{}{
		"AudioID=": audioID,
	}
	parentKey := clouddatastore.NameKey(parentKind, userID, nil)
	q := &model.Query{
		Kind:      kind,
		Parent:    parentKey,
		Limit:     1,
		Cursor:    "",
		Filters:   filters,
		OrderBy:   "",
		Namespace: "",
	}
	keys, _, _, err := r.client.Run(ctx, q)
	if err != nil {
		return nil, err
	}
	var item *entity

	if err := r.client.Get(ctx, keys[0], &item); err != nil {
		return nil, err
	}

	return item.toDomain(), nil
}

func (r *repo) Exists(ctx context.Context, userID string, id string) (bool, error) {
	parentKey := clouddatastore.NameKey(parentKind, userID, nil)
	key := clouddatastore.NameKey(kind, id, parentKey)
	exists, err := r.client.Exists(ctx, key, onlyID(id))
	if err != nil && !errors.Is(err, clouddatastore.ErrNoSuchEntity) {
		return false, errors.WithStack(err)
	}

	return exists, nil
}

func (r *repo) Put(ctx context.Context, userID string, item *model.Feed) error {
	parentKey := clouddatastore.NameKey(parentKind, userID, nil)
	key := clouddatastore.NameKey(kind, string(item.ID()), parentKey)
	if err := r.client.Put(ctx, key, toEntity(item)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repo) PutMulti(ctx context.Context, feeds []*model.Feed) error {
	keys := make([]*clouddatastore.Key, len(feeds))
	entities := make([]*entity, len(feeds))

	for i, item := range feeds {
		parentKey := clouddatastore.NameKey(parentKind, item.UserID, nil)
		keys[i] = clouddatastore.NameKey(kind, string(item.ID()), parentKey)
		entities[i] = toEntity(item)
	}
	err := r.client.PutMulti(ctx, keys, feeds)
	return err
}

func (r *repo) PutTx(tx *clouddatastore.Transaction, userID string, item *model.Feed) error {
	parentKey := clouddatastore.NameKey(parentKind, userID, nil)
	key := clouddatastore.NameKey(kind, string(item.ID()), parentKey)

	if err := r.client.PutTx(tx, key, toEntity(item)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
//
//// TODO: idに型をつけよう。。
func (r *repo) DeleteTx(tx *clouddatastore.Transaction, userID string, id string) error {
	pKey := clouddatastore.NameKey(parentKind, userID, nil)
	key := clouddatastore.NameKey(kind, id, pKey)

	if err := r.client.DeleteTx(tx, key); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, userID string, id string) error {
	pKey := clouddatastore.NameKey(parentKind, userID, nil)
	key := clouddatastore.NameKey(kind, id, pKey)

	if err := r.client.Delete(ctx, key); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
