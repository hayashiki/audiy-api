package ds

import (
	"context"
	"os"
	"testing"
	"time"

	entity2 "github.com/hayashiki/audiy-api/src/domain/entity"
)

func TestFindAudio(t *testing.T) {
	ctx := context.Background()
	dsCli, _ := NewClient(ctx, os.Getenv("GCP_PROJECT"))
	audioRepo := audioRepository{dsCli}
	//playRepo := playRepository{dsCli}
	testAudios(t, &audioRepo)
}

func testAudios(t *testing.T, repo entity2.AudioRepository) {
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

	result, err := repo.Find(ctx, "14145")
	//t.Log(result.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)

	//err = repo.Remove(ctx, result)
	//if err != nil {
	//	t.Fatal(err)
	//}

	if err := repo.Save(ctx, audio); err != nil {
		t.Fatal(err)
	}

	result, err = repo.Find(ctx, "14145")
	if err != nil {
		t.Fatal(err)
	}
	if actual, want := result.ID, audio.ID; actual != want {
		t.Errorf("find result, actual: %v, want: %v", actual, want)
	}
}
