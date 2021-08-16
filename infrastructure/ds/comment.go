package ds

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/datastore"
	"github.com/hayashiki/audiy-api/domain/entity"
	"google.golang.org/api/iterator"
)

type commentRepository struct {
	client *datastore.Client
}

func NewCommentRepository(client *datastore.Client) entity.CommentRepository {
	return &commentRepository{client: client}
}

// GetAll user
func (repo *commentRepository) GetAll(ctx context.Context, userID string, audioID string, cursor string, limit int, sort ...string) ([]*entity.Comment, string, error) {
	//userKey := entity.GetUserKey(userID)
	audioKey := entity.GetAudioKey(audioID)
	query := datastore.NewQuery(entity.CommentKind).Filter("audio_key=", audioKey)

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
	//for _, filter := range opts.Filters {
	//	query = query.Filter(filter.key, filter.value)
	//}
	for _, order := range sort {
		query = query.Order(order)
	}
	it := repo.client.Run(ctx, query)
	entities := make([]*entity.Comment, 0)
	for {
		entity := &entity.Comment{}

		_, err := it.Next(entity)
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Fatalln(err)
			return entities, "", err
		}
		entity.ID = entity.Key.ID
		entities = append(entities, entity)
	}

	nextCursor, err := it.Cursor()
	if err != nil {
		return entities, "", err
	}

	return entities, nextCursor.String(), nil
}

// Save saves comment
func (repo *commentRepository) Save(ctx context.Context, comment *entity.Comment) error {
	// TODO: if exists
	key, err := repo.client.Put(ctx, datastore.IncompleteKey(entity.CommentKind, nil), comment)
	comment.Key = key
	comment.ID = key.ID
	return err
}

// Delete deletes comment
func (repo *commentRepository) Delete(ctx context.Context, commentKey *datastore.Key) error {
	_, err := repo.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		err := repo.client.Delete(ctx, commentKey)
		return err
	})
	return err
}
