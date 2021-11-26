package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/hayashiki/audiy-api/src/domain/entity"
	"github.com/hayashiki/audiy-api/src/infrastructure/gcs"
)

type AudioUsecase interface {
	GetConnection(ctx context.Context, cursor string, limit int, order []string) (*entity.AudioConnection, error)
	Get(ctx context.Context, id string) (*entity.Audio, error)
	CreateAudio(ctx context.Context, input *entity.CreateAudioInput) (*entity.Audio, error)
	UploadAudio(ctx context.Context, input *entity.UploadAudioInput) (*entity.Audio, error)
}

func NewAudioUsecase(
	gcsSvc gcs.Service,
	audioRepo entity.AudioRepository,
	feedRepo entity.FeedRepository,
	userRepo entity.UserRepository,
) AudioUsecase {
	return &audioUsecase{
		gcsSvc:    gcsSvc,
		audioRepo: audioRepo,
		feedRepo:  feedRepo,
		userRepo:  userRepo,
	}
}

type audioUsecase struct {
	gcsSvc    gcs.Service
	audioRepo entity.AudioRepository
	feedRepo  entity.FeedRepository
	userRepo  entity.UserRepository
}

func (u *audioUsecase) CreateAudio(ctx context.Context, input *entity.CreateAudioInput) (*entity.Audio, error) {
	newAudio :=  entity.NewAudio(input.ID, input.Name, input.Length, input.URL, input.Mimetype, time.Now())

	err := u.audioRepo.Save(ctx, newAudio)
	log.Printf("newAudio %+v", newAudio.GetKey())
	if err != nil {
		return nil, fmt.Errorf("fail to create radios record err: %w", err)
	}
	return newAudio, err
}

func (u *audioUsecase) UploadAudio(ctx context.Context, input *entity.UploadAudioInput) (*entity.Audio, error) {
	genID := "TESTID"
	log.Println("description", input.Description)
	b := bytes.Buffer{}
	if _, err := io.Copy(&b, input.File.File); err != nil {
		return nil, err
	}
	if err := u.gcsSvc.Write(ctx, genID, input.File.File); err != nil {
		return nil, err
	}

	// 一旦テスト的にここでとめる
	return nil, nil

	newAudio := entity.NewAudio(genID, input.File.Filename, int(100), "dummy", input.File.ContentType, time.Now())

	err := u.audioRepo.Save(ctx, newAudio)
	log.Printf("newAudio %+v", newAudio.GetKey())
	if err != nil {
		return nil, fmt.Errorf("fail to create radios record err: %w", err)
	}

	users, _ := u.userRepo.GetAll(ctx)
	feeds := make([]*entity.Feed, len(users))
	userIDs := make([]string, len(users))
	newFeed := entity.NewFeed(newAudio.Key.Name, newAudio.PublishedAt)
	newFeed.PublishedAt = newAudio.PublishedAt

	for i, u := range users {
		userIDs[i] = u.ID
		feeds[i] = newFeed
	}
	err = u.feedRepo.SaveAll(ctx, userIDs, feeds)
	return newAudio, err
}

func (u *audioUsecase) GetConnection(ctx context.Context, cursor string, limit int, order []string) (*entity.AudioConnection, error) {
	audios, nextCursor, err := u.audioRepo.FindAll(ctx, nil, cursor, limit, order...)
	if err != nil {
		return nil, err
	}
	audioEdges := make([]*entity.AudioEdge, len(audios))
	for i, a := range audios {
		audioEdges[i] = &entity.AudioEdge{
			Cursor: nextCursor,
			Node:   a,
		}
	}
	return &entity.AudioConnection{
		PageInfo: &entity.PageInfo{
			Cursor:  nextCursor,
			HasMore: len(audios) != 0,
		},
		Edges: audioEdges,
	}, nil
}

func (u *audioUsecase) Get(ctx context.Context, id string) (*entity.Audio, error) {
	return u.audioRepo.Find(ctx, id)
}
