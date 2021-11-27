package usecase

import (
	"context"

	"github.com/hayashiki/audiy-api/src/domain/entity"
)

type MockUserUsecase struct {
	SaveFunc func(input entity.CreateUserInput) (*entity.User, error)
	GetFunc  func(id string) (*entity.User, error)
}

func (m MockUserUsecase) Save(ctx context.Context, input entity.CreateUserInput) (*entity.User, error) {
	return m.SaveFunc(input)
}

func (m MockUserUsecase) Get(ctx context.Context, id string) (*entity.User, error) {
	panic("implement me")
}

func NewMockUserUsecase() UserUsecase {
	return &MockUserUsecase{}
}

func sample() {
	//input := entity.CreateUserInput{
	//	ID:    "11",
	//	Email: "hh@hayashiki.com",
	//}

	userUsecase := MockUserUsecase{}
	userUsecase.SaveFunc = func(input entity.CreateUserInput) (*entity.User, error) {
		return nil, nil
	}
}
