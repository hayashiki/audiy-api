package usecase

import (
	"context"
	"github.com/hayashiki/audiy-api/domain/entity"
)

type PlayUsecase interface {
	Exists(ctx context.Context, userID int64, audioID string) (bool, error)
	Save(ctx context.Context, userID int64, audioID string) (*entity.Play, error)
}

func NewPlayUsecase(audioUserRepo entity.PlayRepository) PlayUsecase {
	return &playUsecase{audioUserRepo: audioUserRepo}
}

type playUsecase struct {
	audioUserRepo entity.PlayRepository
}

func (uc *playUsecase) Exists(ctx context.Context, userID int64, audioID string) (bool, error) {
	return uc.audioUserRepo.Exists(ctx, userID, audioID)
}

func (uc *playUsecase) Save(ctx context.Context, userID int64, audioID string) (*entity.Play, error) {
	au := entity.NewPlay(userID, audioID)
	err := uc.audioUserRepo.Save(ctx, au)
	if err != nil {
		return nil, err
	}
	return au, err
}
