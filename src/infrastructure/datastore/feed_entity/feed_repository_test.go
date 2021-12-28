package feed_entity

import (
	"context"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore/user_entity"
	"log"
	"testing"
)

func TestPutMulti(t *testing.T) {
	ctx := context.Background()
	ds := datastore.New()
	ds2 := datastore.NewDS()
	feedRepo := repo{client: ds2}

	userRepo := user_entity.NewUserRepository(ds)

	user, err := userRepo.Get(ctx, "103843140833205663533")
	if err != nil {
		t.Error(err)
	}

	feeds, _, _, _ := feedRepo.GetAll(ctx, user.ID, nil,"", 200, "CreatedAt")
	for _, f := range feeds {
		log.Println(f)
	}

}
