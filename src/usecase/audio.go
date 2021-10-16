package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	entity2 "github.com/hayashiki/audiy-api/src/domain/entity"
	gcs2 "github.com/hayashiki/audiy-api/src/infrastructure/gcs"
)

type AudioUsecase interface {
	GetConnection(ctx context.Context, cursor string, limit int, order []string) (*entity2.AudioConnection, error)
	Get(ctx context.Context, id string) (*entity2.Audio, error)
	CreateAudio(ctx context.Context, input *entity2.CreateAudioInput) (*entity2.Audio, error)
}

func NewAudioUsecase(
	gcsSvc gcs2.Client,
	audioRepo entity2.AudioRepository,
	feedRepo entity2.FeedRepository,
	userRepo entity2.UserRepository,
) AudioUsecase {
	return &audioUsecase{
		gcsSvc:    gcsSvc,
		audioRepo: audioRepo,
		feedRepo:  feedRepo,
		userRepo:  userRepo,
	}
}

type audioUsecase struct {
	gcsSvc    gcs2.Client
	audioRepo entity2.AudioRepository
	feedRepo  entity2.FeedRepository
	userRepo  entity2.UserRepository
}

func (u *audioUsecase) CreateAudio(ctx context.Context, input *entity2.CreateAudioInput) (*entity2.Audio, error) {
	genID := "TESTID"
	log.Println("description", input.Description)
	b := bytes.Buffer{}
	if _, err := io.Copy(&b, input.File.File); err != nil {
		return nil, err
	}
	if err := u.gcsSvc.Put(ctx, genID, b.Bytes()); err != nil {
		return nil, err
	}

	// 一旦テスト的にここでとめる
	return nil, nil

	newAudio := entity2.NewAudio(genID, input.File.Filename, int(100), "dummy", input.File.ContentType, time.Now())

	err := u.audioRepo.Save(ctx, newAudio)
	log.Printf("newAudio %+v", newAudio.GetKey())
	if err != nil {
		return nil, fmt.Errorf("fail to create radios record err: %w", err)
	}

	users, _ := u.userRepo.GetAll(ctx)
	feeds := make([]*entity2.Feed, len(users))
	userIDs := make([]string, len(users))
	newFeed := entity2.NewFeed(newAudio.Key.Name, newAudio.PublishedAt)
	newFeed.PublishedAt = newAudio.PublishedAt

	for i, u := range users {
		userIDs[i] = u.ID
		feeds[i] = newFeed
	}
	err = u.feedRepo.SaveAll(ctx, userIDs, feeds)
	return newAudio, err
}

func (u *audioUsecase) GetConnection(ctx context.Context, cursor string, limit int, order []string) (*entity2.AudioConnection, error) {
	audios, nextCursor, err := u.audioRepo.FindAll(ctx, nil, cursor, limit, order...)
	if err != nil {
		return nil, err
	}
	audioEdges := make([]*entity2.AudioEdge, len(audios))
	for i, a := range audios {
		audioEdges[i] = &entity2.AudioEdge{
			Cursor: nextCursor,
			Node:   a,
		}
	}
	return &entity2.AudioConnection{
		PageInfo: &entity2.PageInfo{
			Cursor:  nextCursor,
			HasMore: len(audios) != 0,
		},
		Edges: audioEdges,
	}, nil
}

func (u *audioUsecase) Get(ctx context.Context, id string) (*entity2.Audio, error) {
	return u.audioRepo.Find(ctx, id)
}
