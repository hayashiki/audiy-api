package ds

import (
	"context"
	"errors"
	"log"

	entity2 "github.com/hayashiki/audiy-api/src/domain/entity"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

// AudioRepository operates Audio entity
type audioRepository struct {
	client *datastore.Client
}

// NewAudioRepository returns the AudioRepository
func NewAudioRepository(client *datastore.Client) entity2.AudioRepository {
	return &audioRepository{client: client}
}

// Exists exists item
func (repo *audioRepository) Exists(ctx context.Context, id string) bool {
	_, err := repo.Find(ctx, id)
	return err == nil
}

// FindAll finds all radios
func (repo *audioRepository) FindAll(ctx context.Context, filters map[string]interface{}, cursor string, limit int, sort ...string) ([]*entity2.Audio, string, error) {
	query := datastore.NewQuery(entity2.AudioKind)
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
	//query = query.Filter("mimetype=", "audio/mp4")

	for _, order := range sort {
		query = query.Order(order)
	}
	log.Printf("query %+v", query)
	it := repo.client.Run(ctx, query)
	entities := make([]*entity2.Audio, 0)
	for {
		entity := &entity2.Audio{}

		_, err := it.Next(entity)
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return entities, "", err
		}
		entity.ID = entity.Key.Name
		entities = append(entities, entity)
	}

	nextCursor, err := it.Cursor()
	log.Println(nextCursor, err)
	if err != nil {
		return entities, "", err
	}

	return entities, nextCursor.String(), nil
}

// Find finds audio given id
func (repo *audioRepository) Find(ctx context.Context, id string) (*entity2.Audio, error) {
	var audio entity2.Audio
	err := repo.client.Get(ctx, datastore.NameKey(entity2.AudioKind, id, nil), &audio)
	audio.ID = audio.Key.Name
	return &audio, err
}

// Save saves audios
func (repo *audioRepository) Save(ctx context.Context, audio *entity2.Audio) error {
	key, err := repo.client.Put(ctx, datastore.NameKey(entity2.AudioKind, audio.ID, nil), audio)
	audio.Key = key

	return err
}

// Delete saves audios
func (repo *audioRepository) Delete(ctx context.Context, audioKey *datastore.Key) error {
	err := repo.client.Delete(ctx, audioKey)
	return err
}
