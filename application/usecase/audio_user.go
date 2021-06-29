package usecase

import (
	"context"
	"github.com/hayashiki/audiy-api/domain/entity"
)

type AudioUserUsecase interface {
	Save(ctx context.Context, userID int64, audioID string) (*entity.AudioUser, error)
}

func NewAudioUserUsecase(audioUserRepo entity.AudioUserRepository) AudioUserUsecase {
	return &audioUserUsecase{audioUserRepo: audioUserRepo}
}

type audioUserUsecase struct {
	audioUserRepo entity.AudioUserRepository
}

func (uc *audioUserUsecase) Save(ctx context.Context, userID int64, audioID string) (*entity.AudioUser, error) {
	au := entity.NewAudioUser(userID, audioID)
	err := uc.audioUserRepo.Save(ctx, au)
	if err != nil {
		return nil, err
	}
	return au, err
}
