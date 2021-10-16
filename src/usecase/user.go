package usecase

import (
	"context"

	entity2 "github.com/hayashiki/audiy-api/src/domain/entity"
)

type UserUsecase interface {
	Save(context.Context, entity2.CreateUserInput) (*entity2.User, error)
	Get(ctx context.Context, id string) (*entity2.User, error)
}

func NewUserUsecase(
	userRepo entity2.UserRepository,
	audioRepo entity2.AudioRepository,
	feedRepo entity2.FeedRepository,
) UserUsecase {
	return &userUsecase{userRepo: userRepo, audioRepo: audioRepo, feedRepo: feedRepo}
}

type userUsecase struct {
	userRepo  entity2.UserRepository
	audioRepo entity2.AudioRepository
	feedRepo  entity2.FeedRepository
}

func (c *userUsecase) Save(ctx context.Context, input entity2.CreateUserInput) (*entity2.User, error) {
	newUser := entity2.NewUser(input.ID, input.Email, input.Name, input.PhotoURL)
	// exists check
	user, err := c.userRepo.Get(ctx, input.ID)
	if user != nil {
		return user, err
	}
	err = c.userRepo.Save(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// Publishで非同期でもよい
	audios, _, _ := c.audioRepo.FindAll(ctx, nil, "", 1000, "-published_at")
	feeds := make([]*entity2.Feed, len(audios))
	userIDs := make([]string, len(audios))
	for i, a := range audios {
		newFeed := entity2.NewFeed(a.Key.Name, a.PublishedAt)
		feeds[i] = newFeed
		userIDs[i] = newUser.ID
	}
	err = c.feedRepo.SaveAll(ctx, userIDs, feeds)
	return newUser, err
}

func (c *userUsecase) Get(ctx context.Context, id string) (*entity2.User, error) {
	return c.userRepo.Get(ctx, id)
}
