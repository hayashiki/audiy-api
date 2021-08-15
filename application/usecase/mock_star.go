package usecase

import (
	"context"

	"github.com/hayashiki/audiy-api/domain/entity"
)

type MockStarUsecase struct {
	//SaveFunc  func(input entity.CreateUserInput) error
	ExistsFunc func(userID string, audioID string) (bool, error)
}

func (m MockStarUsecase) Delete(ctx context.Context, userID string, audioID string) (*entity.Star, error) {
	panic("implement me")
}

func (m MockStarUsecase) Exists(ctx context.Context, userID string, audioID string) (bool, error) {
	return m.ExistsFunc(userID, audioID)
}

func (m MockStarUsecase) Save(ctx context.Context, userID string, audioID string) (*entity.Star, error) {
	panic("implement me")
}

func NewMockStarUsecase() StarUsecase {
	return &MockStarUsecase{}
}
