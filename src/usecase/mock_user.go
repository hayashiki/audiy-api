package usecase

import (
	"context"

	"github.com/hayashiki/audiy-api/src/domain/model"
)

type MockUserUsecase struct {
	SaveFunc func(input model.CreateUserInput) (*model.User, error)
	GetFunc  func(id string) (*model.User, error)
}

func (m MockUserUsecase) Save(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	return m.SaveFunc(input)
}

func (m MockUserUsecase) Get(ctx context.Context, id string) (*model.User, error) {
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
	userUsecase.SaveFunc = func(input model.CreateUserInput) (*model.User, error) {
		return nil, nil
	}
}
