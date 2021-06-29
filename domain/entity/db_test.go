package entity

import (
	"context"
	"github.com/hayashiki/audiy-api/infrastructure/ds"
	"testing"
	"time"
)

func TestFindAudio(t *testing.T) {
	dsDataSource := ds.Connect()
	repo := NewAudioRepository(dsDataSource)

	testAudios(t, repo)
}

func testAudios(t *testing.T, repo AudioRepository) {
	ctx := context.Background()
	pub := time.Now()
	audio := &Audio{
		ID:          "14145",
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
