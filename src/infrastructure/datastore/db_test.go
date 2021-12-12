package datastore

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/repository"
	"github.com/hayashiki/audiy-api/src/infrastructure/datastore/audio_entity"
	"os"
	"testing"
	"time"

	entity2 "github.com/hayashiki/audiy-api/src/domain/model"
)

func TestFindAudio(t *testing.T) {
	ctx := context.Background()
	dsCli, _ := NewClient(ctx, os.Getenv("GCP_PROJECT"))
	audioRepo := audio_entity.NewAudioRepository(dsCli)
	//playRepo := playRepository{dsCli}
	testAudios(t, &audioRepo)
}

func testAudios(t *testing.T, repo repository.AudioRepository) {
	ctx := context.Background()
	pub := time.Now()
	audio := &entity2.Audio{
		ID: "14145",
		//Name:        "hoge",
		//Length:      100,
		//URL:         "https://example.com",
		//Mimetype:    "mp4",
		PublishedAt: pub,
	}

	result, err := repo.Get(ctx, "14145")
	//t.Log(result.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)

	//err = repo.Remove(ctx, result)
	//if err != nil {
	//	t.Fatal(err)
	//}

	if err := repo.Put(ctx, audio); err != nil {
		t.Fatal(err)
	}

	result, err = repo.Get(ctx, "14145")
	if err != nil {
		t.Fatal(err)
	}
	if actual, want := result.ID, audio.ID; actual != want {
		t.Errorf("find result, actual: %v, want: %v", actual, want)
	}
}
