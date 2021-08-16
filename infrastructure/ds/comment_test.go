package ds

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/hayashiki/audiy-api/domain/entity"
)

func TestCommentSave(t *testing.T) {
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

	commentRepo := commentRepository{dsCli}
	t.Log(user)
	t.Log(audio.ID)
	newComment := entity.NewComment(user.ID, audio.ID, "hogehoge")
	if err := commentRepo.Save(ctx, newComment); err != nil {
		t.Error(err)
	}

	comments, nextCursor, err := commentRepo.GetAll(ctx, "", 2, "id")
	log.Println(comments)
	log.Println(nextCursor)
	log.Println(err)
}
