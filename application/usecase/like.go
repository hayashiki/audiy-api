package usecase

import (
	"context"

	"github.com/hayashiki/audiy-api/domain/entity"
)

type LikeUsecase interface {
	Toggle(ctx context.Context, userID int64, audioID string) (*entity.ToggleLikeResult, error)
	Exists(ctx context.Context, userID int64, audioID string) (bool, error)
	Save(ctx context.Context, userID int64, audioID string) (*entity.Like, error)
	Delete(ctx context.Context, userID int64, audioID string) (*entity.Like, error)
}

func NewLikeUsecase(likeRepo entity.LikeRepository) LikeUsecase {
	return &likeUsecase{likeRepo: likeRepo}
}

type likeUsecase struct {
	likeRepo entity.LikeRepository
}

func (l *likeUsecase) Toggle(ctx context.Context, userID int64, audioID string) (*entity.ToggleLikeResult, error) {
	panic("implement me")
}

func (l *likeUsecase) Exists(ctx context.Context, userID int64, audioID string) (bool, error) {
	panic("implement me")
}

func (l *likeUsecase) Save(ctx context.Context, userID int64, audioID string) (*entity.Like, error) {
	panic("implement me")
}

func (l *likeUsecase) Delete(ctx context.Context, userID int64, audioID string) (*entity.Like, error) {
	panic("implement me")
}
