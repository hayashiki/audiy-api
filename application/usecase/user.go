package usecase

import (
	"context"

	"github.com/hayashiki/audiy-api/domain/entity"
)

type UserUsecase interface {
	Save(context.Context, entity.CreateUserInput) (*entity.User, error)
	Get(ctx context.Context, id string) (*entity.User, error)
}

func NewUserUsecase(
	userRepo entity.UserRepository,
	audioRepo entity.AudioRepository,
	feedRepo entity.FeedRepository,
) UserUsecase {
	return &userUsecase{userRepo: userRepo, audioRepo: audioRepo, feedRepo: feedRepo}
}

type userUsecase struct {
	userRepo  entity.UserRepository
	audioRepo entity.AudioRepository
	feedRepo  entity.FeedRepository
}

func (c *userUsecase) Save(ctx context.Context, input entity.CreateUserInput) (*entity.User, error) {
	newUser := entity.NewUser(input.ID, input.Email)
	err := c.userRepo.Save(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// Publishで非同期でもよい
	audios, _, _ := c.audioRepo.FindAll(ctx, nil, "", 1000, "-published_at")
	feeds := make([]*entity.Feed, len(audios))
	for i, a := range audios {
		newFeed := entity.NewFeed(a.Key.Name, a.PublishedAt)
		feeds[i] = newFeed
	}
	err = c.feedRepo.SaveAll(ctx, []string{newUser.ID}, feeds)
	return newUser, err
}

func (c *userUsecase) Get(ctx context.Context, id string) (*entity.User, error) {
	return c.userRepo.Get(ctx, id)
}
