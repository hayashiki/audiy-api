package user_entity

import (
	"context"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
	"log"
	"os"
	"testing"

	"github.com/hayashiki/audiy-api/src/domain/model"
)

func TestSaveAndGetUser(t *testing.T) {
	log.Println(os.Getenv("GCP_PROJECT"))
	ctx := context.Background()
	ds := datastore.New()
	userRepo := repo{client: ds}


	var id = "111111"
	user := model.NewUser(id, "hayashiki@example.com", "hayashiki", "http://example.com/profile.png")
	if err := userRepo.Put(ctx, user); err != nil {
		t.Fatal(err)
	}

	exists, err := userRepo.Exists(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(exists)

	//if !exists {
	//	t.Fatal("not exists user data")
	//}

	users, err := userRepo.GetAll(ctx)
	for _, u := range users {
		log.Println(u.ID)
		userRepo.Put(ctx, u)
	}
}
