package usecase

import (
	"context"
	"github.com/hayashiki/audiy-api/domain/entity"
)

type LikeUsecase interface {
	Exists(ctx context.Context, userID int64, audioID string) (bool, error)
	Save(ctx context.Context, userID int64, audioID string) (*entity.Like, error)
	Delete(ctx context.Context, userID int64, audioID string) (*entity.Like, error)
}
