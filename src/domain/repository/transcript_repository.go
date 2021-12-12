package repository

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/model"
)

// TranscriptRepository interface
type TranscriptRepository interface {
	GetAll(ctx context.Context) ([]*model.Transcript, error)
	Put(context.Context, *model.Transcript) error
}
