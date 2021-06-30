package entity

import (
	"context"
	"github.com/hayashiki/audiy-api/infrastructure/ds"
	"log"
	"os"
	"testing"
)

func TestAudioUser(t *testing.T) {
	log.Println(os.Getenv("GCP_PROJECT"))
	ctx := context.Background()
	dsDataSource := ds.Connect()
	audioRepo := NewAudioRepository(dsDataSource)
	audio, _ := audioRepo.Find(ctx, "")

	audioUser := NewAudioUser(111111, audio.ID)

	audioUserRepo := audioUserRepository{dsDataSource}
	t.Log(audioUser)
	err := audioUserRepo.Save(ctx, audioUser)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAnotherAudioUser(t *testing.T) {
	log.Println(os.Getenv("GCP_PROJECT"))
	ctx := context.Background()
	dsDataSource := ds.Connect()
	audioRepo := NewAudioRepository(dsDataSource)
	audio, _ := audioRepo.Find(ctx, "")

	audioUser := NewAudioUser(111111, audio.ID)

	audioUserRepo := audioUserRepository{dsDataSource}
	t.Log(audioUser)
	err := audioUserRepo.Save(ctx, audioUser)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExistsAudioUser(t *testing.T) {
	log.Println(os.Getenv("GCP_PROJECT"))
	ctx := context.Background()
	dsDataSource := ds.Connect()
	audioRepo := NewAudioRepository(dsDataSource)
	audio, _ := audioRepo.Find(ctx, "")
	audioUserRepo := audioUserRepository{dsDataSource}
	_, err := audioUserRepo.Exists(ctx, 111111, audio.ID)
	if err != nil {
		t.Fatal(err)
	}
}


