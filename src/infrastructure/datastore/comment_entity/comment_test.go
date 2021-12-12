package comment_entity

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore/audio_entity"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore/user_entity"
	"log"
	"testing"
)

var ds datastore.Client

func init()  {
	ds = datastore.New()
}

func TestGetAll(t *testing.T) {
	//log.Println(os.Getenv("GCP_PROJECT"))
	ds := datastore.New()
	r := repo{client: ds}

	ctx := context.Background()

	userRepo := user_entity.NewUserRepository(ds)
	user, err := userRepo.Get(ctx, "103843140833205663533")
	if err != nil {
		t.Error(err)
	}
	t.Log(user)

	audioRepo := audio_entity.NewAudioRepository(ds)
	if err != nil {
		t.Error(err)
	}
	audio, err := audioRepo.Get(ctx, "F02D2M3L1C7")
	if err != nil {
		t.Error(err)
	}
	t.Log(audio.PublishedAt)

	comments, nextCursor, _, err := r.GetAllByAudio(
		ctx, audio.ID, "", 100, "CreatedAt",
	)

	log.Println(comments)
	log.Println(nextCursor)
	log.Println(err)

	for _, cm := range comments {
		log.Println(cm.Body)
	}
}

func TestPut(t *testing.T) {
	ctx := context.Background()
	r := repo{client: ds}

	userRepo := user_entity.NewUserRepository(ds)
	user, err := userRepo.Get(ctx, "103843140833205663533")
	if err != nil {
		t.Error(err)
	}
	t.Log(user)

	audioRepo := audio_entity.NewAudioRepository(ds)
	if err != nil {
		t.Error(err)
	}
	audio, err := audioRepo.Get(ctx, "F02D2M3L1C7")
	if err != nil {
		t.Error(err)
	}

	newComment := model.NewComment(user.ID, audio.ID, "TestBody")

	err = r.Put(ctx, newComment)
	if err != nil {
		t.Error(err)
	}
}
