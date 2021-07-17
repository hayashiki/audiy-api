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

func NewStarUsecase(starRepo entity.StarRepository) StarUsecase {
	return &starUsecase{starRepo: starRepo}
}

type starUsecase struct {
	starRepo entity.StarRepository
}

func (s *starUsecase) Exists(ctx context.Context, userID int64, audioID string) (bool, error) {
	return s.starRepo.Exists(ctx, userID, audioID)
}

func (s *starUsecase) Save(ctx context.Context, userID int64, audioID string) (*entity.Star, error) {
	newsStar := entity.NewStar(userID, audioID)
	err := s.starRepo.Save(ctx, newsStar)
	return newsStar, err
}

func (s *starUsecase) Delete(ctx context.Context, userID int64, audioID string) (*entity.Star, error) {
	star, err := s.starRepo.FindByRel(ctx, userID, audioID)
	err = s.starRepo.Delete(ctx, star.ID)
	if err != nil {
		return nil, err
	}
	return star, err
}
