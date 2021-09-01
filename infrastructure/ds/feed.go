package ds

import (
	"context"
	"errors"
	"fmt"
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
func (repo *feedRepository) FindAll(ctx context.Context, userID string, filters map[string]interface{}, cursor string, limit int, sort ...string) ([]*entity.Feed, string, bool, error) {
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
		query = query.Limit(limit + 1)
	}
	for key, val := range filters {
		log.Println(key, val)
		query = query.Filter(key+"=", val)
	}

	for _, order := range sort {
		query = query.Order(order)
	}
	log.Printf("query %+v", query)
	it := repo.client.Run(ctx, query)
	entities := make([]*entity.Feed, 0, limit)
	count := 0
	hasMore := false
	var nextCursor datastore.Cursor
	for {
		entity := &entity.Feed{}

		_, err := it.Next(entity)
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Printf("err %+v", err)
			return entities, "", hasMore, err
		}
		count++
		if limit < count {
			hasMore = true
			break
		}

		entity.ID = entity.Key.ID
		entities = append(entities, entity)

		if limit == count {
			nextCursor, err = it.Cursor()
			if err != nil {
				return entities, "", hasMore, err
			}
		}
	}
	log.Printf("entities %+v", len(entities))
	return entities, nextCursor.String(), hasMore, nil
}

// Find finds Feed given id
func (repo *feedRepository) Find(ctx context.Context, id int64, userID string) (*entity.Feed, error) {
	var feed entity.Feed
	err := repo.client.Get(ctx, repo.getKey(id, userID), &feed)
	feed.ID = feed.Key.ID
	return &feed, err
}

func (repo *feedRepository) FindByAudio(ctx context.Context, userID string, audioID string) (*entity.Feed, error) {
	userKey := entity.GetUserKey(userID)
	audioKey := entity.GetAudioKey(audioID)
	q := datastore.NewQuery(entity.FeedKind).Ancestor(userKey).KeysOnly().Filter("audio_key =", audioKey).Limit(1)

	keys, err := repo.client.GetAll(context.Background(), q, nil)

	if err != nil {
		return nil, fmt.Errorf("not found feed keys %w", err)
	}

	var feed entity.Feed

	if err := repo.client.Get(ctx, keys[0], &feed); err != nil {
		return nil, fmt.Errorf("not found feed %w", err)
	}

	feed.ID = keys[0].ID

	return &feed, nil
}

// Save saves Feeds
func (repo *feedRepository) Save(ctx context.Context, userID string, feed *entity.Feed) error {
	if feed.Key != nil {
		_, err := repo.client.Put(ctx, feed.Key, feed)
		return err
	}

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
