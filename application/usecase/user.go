package usecase

import (
	"context"
	"github.com/hayashiki/audiy-api/domain/entity"
	"strconv"
)

type UserUsecase interface {
	Save(context.Context, entity.CreateUserInput) error
	Get(ctx context.Context, id int64) (*entity.User, error)
}

func NewUserUsecase(userRepo entity.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

type userUsecase struct {
	userRepo entity.UserRepository
}

func (c *userUsecase) Save(ctx context.Context, input entity.CreateUserInput) error {
	id, _ := strconv.Atoi(input.ID)
	newUser := entity.NewUser(int64(id), input.Email)
	err := c.userRepo.Save(ctx, newUser)
	if err != nil {
		return err
	}
	return err
}

func (c *userUsecase) Get(ctx context.Context, id int64) (*entity.User, error) {
	return c.userRepo.Get(ctx, id)
}
