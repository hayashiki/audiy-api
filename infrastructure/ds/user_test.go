package ds

import (
	"context"
	"github.com/hayashiki/audiy-api/domain/entity"
	"log"
	"os"
	"testing"
)

func TestSaveAndGetUser(t *testing.T) {
	log.Println(os.Getenv("GCP_PROJECT"))
	ctx := context.Background()
	dsCli, _ := NewClient(ctx, os.Getenv("GCP_PROJECT"))
	audioUserRepo := userRepository{dsCli}
	audioRepo := audioRepository{dsCli}
	playRepo := playRepository{dsCli}

	var id int64 = 111111
	user := entity.NewUser(id, "hayashiki@example.com")
	err := audioUserRepo.Save(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	userKey := entity.GetUserKey(id)
	exists, err := audioUserRepo.Exists(ctx, userKey)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("not exists user data")
	}
	audio, _ := audioRepo.Find(ctx, "F023GTZRRU2")
	exists, err = playRepo.Exists(ctx, user.ID, audio.ID)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(exists)
	if !exists {
		newPlay := entity.NewPlay(user.ID, audio.ID)
		if err := playRepo.Save(ctx, newPlay); err != nil {
			t.Fatal(err)
		}
	}
}
