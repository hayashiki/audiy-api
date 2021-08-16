package ds

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/hayashiki/audiy-api/domain/entity"
)

func TestStarRepository_Exists(t *testing.T) {
	log.Println(os.Getenv("GCP_PROJECT"))
	ctx := context.Background()
	dsCli, _ := NewClient(ctx, os.Getenv("GCP_PROJECT"))

	var id string = "111111"
	userRepo := userRepository{dsCli}
	user, err := userRepo.Get(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	audioRepo := audioRepository{dsCli}
	audio, err := audioRepo.Find(ctx, "F023GTZRRU2")
	if err != nil {
		t.Fatal(err)
	}

	starRepo := starRepository{client: dsCli}

	star := entity.NewStar(user.ID, audio.ID)
	err = starRepo.Save(ctx, star)
	if err != nil {
		t.Error(err)
		return
	}

	exists, err := starRepo.Exists(ctx, user.ID, audio.ID)
	if err != nil {
		t.Error(err)
	}
	if got, want := true, exists; got != want {
		t.Errorf("Exists got %v, want %v", got, want)
	}
	err = starRepo.Delete(ctx, star.ID)
	if err != nil {
		t.Error(err)
		return
	}
}
