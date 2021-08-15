package usecase

import (
	"context"

	"github.com/hayashiki/audiy-api/domain/entity"
)

type MockLikeUsecase struct {
	//SaveFunc  func(input entity.CreateUserInput) error
	ExistsFunc func(userID string, audioID string) (bool, error)
}

func (m MockLikeUsecase) Delete(ctx context.Context, userID string, audioID string) (*entity.Like, error) {
	panic("implement me")
}

func (m MockLikeUsecase) Exists(ctx context.Context, userID string, audioID string) (bool, error) {
	return m.ExistsFunc(userID, audioID)
}

func (m MockLikeUsecase) Save(ctx context.Context, userID string, audioID string) (*entity.Like, error) {
	panic("implement me")
}

func NewMockLikeUsecase() LikeUsecase {
	return &MockLikeUsecase{}
}
