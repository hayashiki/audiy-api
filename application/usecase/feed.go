package usecase

import (
	"context"
	"log"

	"github.com/hayashiki/audiy-api/domain/entity"
)

type FeedUsecase interface {
	GetConnection(ctx context.Context, userID string, cursor string, limit int, order []string) (*entity.FeedConnection, error)
	Get(ctx context.Context, id int64, userID string) (*entity.Feed, error)
	Put(ctx context.Context, userID string, audioID string, event entity.FeedEvent) (*entity.Feed, error)
}

func NewFeedUsecase(feedRepo entity.FeedRepository) FeedUsecase {
	return &feedUsecase{feedRepo: feedRepo}
}

type feedUsecase struct {
	feedRepo entity.FeedRepository
}

func (u *feedUsecase) GetConnection(ctx context.Context, userID string, cursor string, limit int, order []string) (*entity.FeedConnection, error) {
	feeds, nextCursor, err := u.feedRepo.FindAll(ctx, userID, nil, cursor, limit, order...)
	if err != nil {
		return nil, err
	}
	feedEdges := make([]*entity.FeedEdge, len(feeds))
	for i, a := range feeds {
		feedEdges[i] = &entity.FeedEdge{
			Cursor: nextCursor,
			Node:   a,
		}
	}
	return &entity.FeedConnection{
		PageInfo: &entity.PageInfo{
			Cursor:  nextCursor,
			HasMore: len(feeds) != 0,
		},
		Edges: feedEdges,
	}, nil
}

func (u *feedUsecase) Get(ctx context.Context, id int64, userID string) (*entity.Feed, error) {
	return u.feedRepo.Find(ctx, id, userID)
}

func (u *feedUsecase) Put(ctx context.Context, userID string, audioID string, event entity.FeedEvent) (*entity.Feed, error) {
	feed, err := u.feedRepo.FindByAudio(ctx, userID, audioID)
	if err != nil {
		return nil, err
	}
	switch event {
	case entity.FeedEventPlayed:
		feed.Played = true
	case entity.FeedEventStared:
		feed.Stared = true
	case entity.FeedEventUnstared:
		feed.Stared = false
	case entity.FeedEventLiked:
		feed.Liked = true
	case entity.FeedEventUnliked:
		feed.Liked = false
	default:
		log.Printf("invalid event %s", event)
	}
	err = u.feedRepo.Save(ctx, userID, feed)
	return feed, nil
}
