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

func NewLikeUsecase(likeRepo entity.LikeRepository) LikeUsecase {
	return &likeUsecase{likeRepo: likeRepo}
}

type likeUsecase struct {
	likeRepo entity.LikeRepository
}

func (l *likeUsecase) Exists(ctx context.Context, userID int64, audioID string) (bool, error) {
	panic("implement me")
}

func (l *likeUsecase) Save(ctx context.Context, userID int64, audioID string) (*entity.Like, error) {
	newLike := entity.NewLike(userID, audioID)
	err := l.likeRepo.Save(ctx, newLike)
	return newLike, err
}

func (l *likeUsecase) Delete(ctx context.Context, userID int64, audioID string) (*entity.Like, error) {
	like, err := l.likeRepo.FindByRel(ctx, userID, audioID)
	if err != nil {
		return nil, err
	}
	err = l.likeRepo.Delete(ctx, like.ID)
	return like, err
}
