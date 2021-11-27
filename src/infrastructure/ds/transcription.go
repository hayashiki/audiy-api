package ds

import (
	"context"
	"log"

	"cloud.google.com/go/datastore"
	"github.com/hayashiki/audiy-api/src/domain/entity"
)

type transcriptRepository struct {
	client *datastore.Client
}

func NewTranscriptRepository(client *datastore.Client) entity.TranscriptRepository {
	return &transcriptRepository{client: client}
}

// Save saves transcription
func (repo *transcriptRepository) Save(ctx context.Context, transcript *entity.Transcript) error {
	log.Println("db save", transcript)
	// TODO: if exists
	key, err := repo.client.Put(ctx, datastore.IncompleteKey(entity.TranscriptKind, nil), transcript)
	if err != nil {
		log.Println("db err", err)
		return err
	}

	transcript.Key = key
	transcript.ID = key.ID
	return err
}
