package ds

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/hayashiki/audiy-api/domain/entity"
)

func TestSaveAndGetAudio(t *testing.T) {
	log.Println(os.Getenv("GCP_PROJECT"))
	ctx := context.Background()
	dsCli, _ := NewClient(ctx, os.Getenv("GCP_PROJECT"))
	userRepo := userRepository{dsCli}
	audioRepo := audioRepository{dsCli}
	//playRepo := playRepository{dsCli}

	user, err := userRepo.Get(ctx, "103843140833205663533")
	if err != nil {
		t.Fatal(err)
	}
	exists, err := userRepo.Exists(ctx, user.Key)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("not exists user data")
	}
	newTestAudio1 := entity.NewAudio("dummy1", "dummy1", 100, "url", "audio/mp4", time.Now())
	err = audioRepo.Save(ctx, newTestAudio1)
	if err != nil {
		t.Fatal(err)
	}
	newTestAudio2 := entity.NewAudio("dummy2", "dummy2", 100, "url", "audio/mp4", time.Now())
	err = audioRepo.Save(ctx, newTestAudio2)
	if err != nil {
		t.Fatal(err)
	}
	//newPlay1 := entity.NewPlay(user.ID, newTestAudio1.ID)
	//if err := playRepo.Save(ctx, newPlay1); err != nil {
	//	t.Fatal(err)
	//}
	filters := map[string]interface{}{
		"played_users": user.Key.Name,
	}
	audios, _, err := audioRepo.FindAll(ctx, filters, "", 1000, "-published_at")

	//for _, a := range audios {
	//	a.PlayedUsers = []string{user.Key.Name}
	//	audioRepo.Save(ctx, a)
	//}

	log.Println(len(audios))
}
