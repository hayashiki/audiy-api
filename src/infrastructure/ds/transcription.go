package ds

import (
	"context"
	
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
func (repo *transcriptRepository) Save(ctx context.Context, transcription *entity.Transcript) error {
	// TODO: if exists
	key, err := repo.client.Put(ctx, datastore.IncompleteKey(entity.TranscriptKind, nil), transcription)
	transcription.Key = key
	transcription.ID = key.ID
	return err
}
