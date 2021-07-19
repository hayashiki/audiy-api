package usecase

import (
	"context"

	"github.com/hayashiki/audiy-api/domain/entity"
)

type UserUsecase interface {
	Save(context.Context, entity.CreateUserInput) (*entity.User, error)
	Get(ctx context.Context, id string) (*entity.User, error)
}

func NewUserUsecase(userRepo entity.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

type userUsecase struct {
	userRepo entity.UserRepository
}

func (c *userUsecase) Save(ctx context.Context, input entity.CreateUserInput) (*entity.User, error) {
	newUser := entity.NewUser(input.ID, input.Email)
	err := c.userRepo.Save(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return newUser, err
}

func (c *userUsecase) Get(ctx context.Context, id string) (*entity.User, error) {
	return c.userRepo.Get(ctx, id)
}
