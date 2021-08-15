package usecase

import (
	"context"

	"github.com/hayashiki/audiy-api/domain/entity"
)

type MockPlayUsecase struct {
	//SaveFunc  func(input entity.CreateUserInput) error
	ExistsFunc func(userID string, audioID string) (bool, error)
}

func (m MockPlayUsecase) Exists(ctx context.Context, userID string, audioID string) (bool, error) {
	return m.ExistsFunc(userID, audioID)
}

func (m MockPlayUsecase) Save(ctx context.Context, userID string, audioID string) (*entity.Play, error) {
	panic("implement me")
}

func NewMockPlayUsecase() PlayUsecase {
	return &MockPlayUsecase{}
}
