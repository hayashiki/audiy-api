package usecase

import (
	"context"
	"errors"
	"github.com/hayashiki/audiy-api/src/domain/repository"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
	"go.mercari.io/datastore/boom"
	"log"

	"github.com/hayashiki/audiy-api/src/domain/model"
)

type UserUsecase interface {
	Save(context.Context, model.CreateUserInput) (*model.User, error)
	Get(ctx context.Context, id string) (*model.User, error)
}

func NewUserUsecase(
	transactor datastore.Transactor,
	userRepo repository.UserRepository,
	audioRepo repository.AudioRepository,
	feedRepo repository.FeedRepository,
) UserUsecase {
	return &userUsecase{transactor: transactor, userRepo: userRepo, audioRepo: audioRepo, feedRepo: feedRepo}
}

type userUsecase struct {
	transactor datastore.Transactor
	userRepo  repository.UserRepository
	audioRepo repository.AudioRepository
	feedRepo  repository.FeedRepository
}

func (c *userUsecase) Save(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	newUser := model.NewUser(input.ID, input.Email, input.Name, input.PhotoURL)

	if err := c.transactor.RunInTransaction(ctx, func(tx *boom.Transaction) error {
		exists, err := c.userRepo.Exists(ctx, input.ID)
		if errors.Is(err, datastore.ErrNoSuchEntity) && err != nil {
			return err
		}
		if exists {
			return errors.New("Already exists...")
		}
		err = c.userRepo.PutTx(tx, newUser)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Println("tx err:", err)
	//	TODO: rollback
		return nil, err
	}

	// Publishで非同期でもよい
	audios, _, _, _ := c.audioRepo.GetAll(ctx, "", 1000, "-PublishedAt")
	feeds := make([]*model.Feed, len(audios))
	// TODO: 本当に不要かみきわめる
	//userIDs := make([]string, len(audios))
	for i, a := range audios {
		newFeed := model.NewFeed(a.ID, newUser.ID, a.PublishedAt)
		feeds[i] = newFeed
		//userIDs[i] = newUser.ID
	}
	err := c.feedRepo.PutMulti(ctx, feeds)
	return newUser, err
}

func (c *userUsecase) Get(ctx context.Context, id string) (*model.User, error) {
	return c.userRepo.Get(ctx, id)
}
