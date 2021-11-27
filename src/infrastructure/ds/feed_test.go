package ds

import (
	"context"
	"testing"
	"time"

	"github.com/hayashiki/audiy-api/src/config"

	"github.com/hayashiki/audiy-api/src/domain/entity"
)

func TestSaveAndGetFeed(t *testing.T) {
	ctx := context.Background()
	dsCli, _ := NewClient(ctx, config.GetProject())
	userRepo := userRepository{dsCli}
	audioRepo := audioRepository{dsCli}
	feedRepo := feedRepository{dsCli}

	audio, err := audioRepo.Find(ctx, "F02A9HFN9AR")
	if err != nil {
		t.Error(err)
	}
	user, err := userRepo.Get(ctx, "103843140833205663533")
	if err != nil {
		t.Error(err)
	}

	newFeed := entity.NewFeed(
		audio.ID,
		time.Now(),
	)

	err = feedRepo.Save(ctx, user.ID, newFeed)
}

func TestPutMulti(t *testing.T) {
	ctx := context.Background()
	dsCli, _ := NewClient(ctx, config.GetProject())
	userRepo := userRepository{dsCli}
	audioRepo := audioRepository{dsCli}
	feedRepo := feedRepository{dsCli}
	audios, _, _ := audioRepo.FindAll(ctx, nil, "", 9, "-published_at")
	users, _ := userRepo.GetAll(ctx)
	feeds := make([]*entity.Feed, len(audios))
	userIDs := make([]string, len(audios))
	for _, u := range users {
		for ii, a := range audios {
			newFeed := entity.NewFeed(a.Key.Name, a.PublishedAt)
			feeds[ii] = newFeed
			userIDs[ii] = u.ID
		}
	}
	feedRepo.SaveAll(ctx, userIDs, feeds)
}

func TestFindAll(t *testing.T) {
	ctx := context.Background()
	dsCli, _ := NewClient(ctx, config.GetProject())
	feedRepo := feedRepository{dsCli}
	feedRepo.FindAll(ctx, "userID", nil, "", 10)
}
