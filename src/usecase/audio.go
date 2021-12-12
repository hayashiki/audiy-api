package usecase

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hayashiki/audiy-api/src/domain/repository"
	"io"
	"log"
	"time"

	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/infrastructure/gcs"
)

type AudioUsecase interface {
	GetConnection(ctx context.Context, cursor string, limit int, orderBy string) (*model.AudioConnection, error)
	Get(ctx context.Context, id string) (*model.Audio, error)
	CreateAudio(ctx context.Context, input *model.CreateAudioInput) (*model.Audio, error)
	UploadAudio(ctx context.Context, input *model.UploadAudioInput) (*model.Audio, error)
}

func NewAudioUsecase(
	gcsSvc gcs.Service,
	audioRepo repository.AudioRepository,
	feedRepo repository.FeedRepository,
	userRepo repository.UserRepository,
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
	audioRepo repository.AudioRepository
	feedRepo  repository.FeedRepository
	userRepo  repository.UserRepository
}

func (u *audioUsecase) CreateAudio(ctx context.Context, input *model.CreateAudioInput) (*model.Audio, error) {
	newAudio :=  model.NewAudio(input.ID, input.Name, input.Length, input.URL, input.Mimetype, time.Now())

	err := u.audioRepo.Put(ctx, newAudio)
	log.Printf("newAudio %+v", newAudio.GetKey())
	if err != nil {
		return nil, fmt.Errorf("fail to create radios record err: %w", err)
	}
	return newAudio, err
}

func (u *audioUsecase) UploadAudio(ctx context.Context, input *model.UploadAudioInput) (*model.Audio, error) {
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

	newAudio := model.NewAudio(genID, input.File.Filename, float64(100), "dummy", input.File.ContentType, time.Now())

	err := u.audioRepo.Put(ctx, newAudio)
	if err != nil {
		return nil, fmt.Errorf("fail to create radios record err: %w", err)
	}

	users, _ := u.userRepo.GetAll(ctx)
	feeds := make([]*model.Feed, len(users))

	for i, u := range users {
		newFeed := model.NewFeed(newAudio.ID, u.ID, newAudio.PublishedAt)
		// TODO: setter
		newFeed.PublishedAt = newAudio.PublishedAt
		feeds[i] = newFeed
	}
	err = u.feedRepo.PutMulti(ctx, feeds)
	return newAudio, err
}

func (u *audioUsecase) GetConnection(ctx context.Context, cursor string, limit int, orderBy string) (*model.AudioConnection, error) {
	audios, nextCursor, hasMore, err := u.audioRepo.GetAll(ctx, cursor, limit, orderBy)
	if err != nil {
		return nil, err
	}
	audioEdges := make([]*model.AudioEdge, len(audios))
	for i, a := range audios {
		audioEdges[i] = &model.AudioEdge{
			Cursor: nextCursor,
			Node:   a,
		}
	}
	return &model.AudioConnection{
		PageInfo: &model.PageInfo{
			Cursor:  nextCursor,
			HasMore: hasMore,
		},
		Edges: audioEdges,
	}, nil
}

func (u *audioUsecase) Get(ctx context.Context, id string) (*model.Audio, error) {
	return u.audioRepo.Get(ctx, id)
}
