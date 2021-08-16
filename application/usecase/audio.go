package usecase

import (
	"context"

	"github.com/hayashiki/audiy-api/domain/entity"
)

type AudioUsecase interface {
	GetConnection(ctx context.Context, cursor string, limit int, order []string) (*entity.AudioConnection, error)
	Get(ctx context.Context, id string) (*entity.Audio, error)
}

func NewAudioUsecase(audioRepo entity.AudioRepository) AudioUsecase {
	return &audioUsecase{audioRepo: audioRepo}
}

type audioUsecase struct {
	audioRepo entity.AudioRepository
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
