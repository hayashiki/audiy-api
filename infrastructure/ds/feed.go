package ds

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/datastore"
	"github.com/hayashiki/audiy-api/domain/entity"
	"google.golang.org/api/iterator"
)

// FeedRepository operates Feed entity
type feedRepository struct {
	client *datastore.Client
}

// NewFeedRepository returns the FeedRepository
func NewFeedRepository(client *datastore.Client) entity.FeedRepository {
	return &feedRepository{client: client}
}

func (repo *feedRepository) key(userID string) *datastore.Key {
	return datastore.IncompleteKey(entity.FeedKind, repo.parentKey(userID))
}

func (repo *feedRepository) getKey(id int64, userID string) *datastore.Key {
	return datastore.IDKey(entity.FeedKind, id, repo.parentKey(userID))
}

func (repo *feedRepository) parentKey(userID string) *datastore.Key {
	return datastore.NameKey(entity.UserKind, userID, nil)
}

// Exists exists item
func (repo *feedRepository) Exists(ctx context.Context, id int64, userID string) bool {
	_, err := repo.Find(ctx, id, userID)
	return err == nil
}

// FindAll finds all Feeds
func (repo *feedRepository) FindAll(ctx context.Context, userID string, filters map[string]interface{}, cursor string, limit int, sort ...string) ([]*entity.Feed, string, error) {
	userKey := entity.GetUserKey(userID)
	query := datastore.NewQuery(entity.FeedKind).Ancestor(userKey)
	if cursor != "" {
		dsCursor, err := datastore.DecodeCursor(cursor)
		if err != nil {
			//TODO
			log.Printf("failed to decode %v", err)
		}
		query = query.Start(dsCursor)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	for key, val := range filters {
		log.Println(key, val)
		query = query.Filter(key+"=", val)
	}
	//query = query.Filter("mimetype=", "Feed/mp4")

	for _, order := range sort {
		query = query.Order(order)
	}
	log.Printf("query %+v", query)
	it := repo.client.Run(ctx, query)
	entities := make([]*entity.Feed, 0)
	for {
		entity := &entity.Feed{}

		_, err := it.Next(entity)
		log.Printf("entity %+v", entity)
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return entities, "", err
		}
		entity.ID = entity.Key.ID
		entities = append(entities, entity)
	}
	log.Printf("entities %+v", entities)

	nextCursor, err := it.Cursor()
	if err != nil {
		return entities, "", err
	}

	return entities, nextCursor.String(), nil
}

// Find finds Feed given id
func (repo *feedRepository) Find(ctx context.Context, id int64, userID string) (*entity.Feed, error) {
	var feed entity.Feed
	err := repo.client.Get(ctx, repo.getKey(id, userID), &feed)
	feed.ID = feed.Key.ID
	return &feed, err
}

// Save saves Feeds
func (repo *feedRepository) Save(ctx context.Context, userID string, feed *entity.Feed) error {
	key, err := repo.client.Put(ctx, repo.key(userID), feed)
	feed.Key = key

	return err
}

// Save saves Feeds
func (repo *feedRepository) SaveAll(ctx context.Context, userIDs []string, feeds []*entity.Feed) error {
	keys := make([]*datastore.Key, len(userIDs))
	for i, u := range userIDs {
		keys[i] = repo.key(u)
	}

	keys, err := repo.client.PutMulti(ctx, keys, feeds)
	log.Println(keys, err)
	//feed.Key = key

	return err
}

// Delete saves Feeds
func (repo *feedRepository) Delete(ctx context.Context, FeedKey *datastore.Key) error {
	err := repo.client.Delete(ctx, FeedKey)
	return err
}
