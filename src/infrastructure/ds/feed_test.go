package ds

import (
	"context"
	"log"
	"testing"

	config2 "github.com/hayashiki/audiy-api/src/config"

	entity2 "github.com/hayashiki/audiy-api/src/domain/entity"
)

func TestSaveAndGetFeed(t *testing.T) {
	ctx := context.Background()
	dsCli, _ := NewClient(ctx, config2.GetProject())
	userRepo := userRepository{dsCli}
	audioRepo := audioRepository{dsCli}
	feedRepo := feedRepository{dsCli}

	audio, err := audioRepo.Find(ctx, "F02A9HFN9AR")
	user, err := userRepo.Get(ctx, "103843140833205663533")

	log.Println(audio, err)
	log.Println(user, err)

	newFeed := entity2.NewFeed(
		audio.ID,
	)

	err = feedRepo.Save(ctx, user.ID, newFeed)
	log.Println(err)
}

func TestPutMulti(t *testing.T) {
	ctx := context.Background()
	dsCli, _ := NewClient(ctx, config2.GetProject())
	userRepo := userRepository{dsCli}
	audioRepo := audioRepository{dsCli}
	feedRepo := feedRepository{dsCli}
	audios, _, _ := audioRepo.FindAll(ctx, nil, "", 9, "-published_at")
	users, _ := userRepo.GetAll(ctx)
	feeds := make([]*entity2.Feed, len(audios))
	userIDs := make([]string, len(audios))
	for _, u := range users {
		for ii, a := range audios {
			newFeed := entity2.NewFeed(a.Key.Name, a.PublishedAt)
			feeds[ii] = newFeed
			userIDs[ii] = u.ID
		}
	}
	feedRepo.SaveAll(ctx, userIDs, feeds)
}

func TestFindAll(t *testing.T) {
	ctx := context.Background()
	dsCli, _ := NewClient(ctx, config2.GetProject())
	feedRepo := feedRepository{dsCli}
	feedRepo.FindAll(ctx, "userID", nil, "", 10)
}
