package transcript_entity

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/domain/repository"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
	"github.com/pkg/errors"
	"go.mercari.io/datastore/boom"
	"log"
)

type repo struct {
	client datastore.Client
}

func NewTranscriptRepository(client datastore.Client) repository.TranscriptRepository {
	return &repo{
		client: client,
	}
}

func (r *repo) GetAll(ctx context.Context) ([]*model.Transcript, error) {
	var entities []*entity
	if err := r.client.GetAll(ctx, kind, &entities); err != nil {
		return nil, err
	}
	transcripts := make([]*model.Transcript, len(entities))
	for i, e := range entities {
		transcripts[i] = e.toDomain()
	}
	return transcripts, nil
}

func (r *repo) Put(ctx context.Context, item *model.Transcript) error {
	if err := r.client.Put(ctx, toEntity(item)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repo) PutTx(tx *boom.Transaction, item *model.Transcript) error {
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


