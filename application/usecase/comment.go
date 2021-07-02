package usecase

import (
	"context"
	"github.com/hayashiki/audiy-api/domain/entity"
)

type CommentUsecase interface {
	Save(ctx context.Context, userID int64, audioID string) (*entity.Comment, error)
	FindAll() error
}
