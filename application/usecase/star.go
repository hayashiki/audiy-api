package usecase

import (
	"context"
	"github.com/hayashiki/audiy-api/domain/entity"
)

type StarUsecase interface {
	Exists(ctx context.Context, userID int64, audioID string) (bool, error)
	Save(ctx context.Context, userID int64, audioID string) (*entity.Star, error)
	Delete(ctx context.Context, userID int64, audioID string) (*entity.Star, error)
}
