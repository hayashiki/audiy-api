package audio_entity

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/hayashiki/audiy-api/src/infrastructure/datastore"
)

func TestSaveAndGetAudio(t *testing.T) {
	ctx := context.Background()
	ds := datastore.New()
	audioRepo := repo{client: ds}

	//exists, err := userRepo.Exists(ctx, user.Key)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if !exists {
	//	t.Fatal("not exists user data")
	//}
	//newTestAudio1 := model.NewAudio("dummy1", "dummy1", 100, "url", "audio/mp4", time.Now())
	//err = audioRepo.Put(ctx, newTestAudio1)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//newTestAudio2 := model.NewAudio("dummy2", "dummy2", 100, "url", "audio/mp4", time.Now())
	//err = audioRepo.Put(ctx, newTestAudio2)
	//if err != nil {
	//	t.Fatal(err)
	//}

	audios, _, _, _ := audioRepo.GetAll(ctx, "", 100, "-PublishedAt")
	for _, a := range audios {
		log.Println(a.Name)
		if err := audioRepo.Put(ctx, a); err != nil {
			time.Sleep(1 * time.Second)
			log.Println(err)
		}
	}
	//log.Println(len(audios))
}
