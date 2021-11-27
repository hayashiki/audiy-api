package ds

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/hayashiki/audiy-api/src/domain/entity"
)

func TestSaveAndGetUser(t *testing.T) {
	log.Println(os.Getenv("GCP_PROJECT"))
	ctx := context.Background()
	dsCli, _ := NewClient(ctx, os.Getenv("GCP_PROJECT"))
	userRepo := userRepository{dsCli}

	var id = "111111"
	user := entity.NewUser(id, "hayashiki@example.com", "hayashiki", "http://example.com/profile.png")
	err := userRepo.Save(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	userKey := entity.GetUserKey(id)
	exists, err := userRepo.Exists(ctx, userKey)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("not exists user data")
	}
}
