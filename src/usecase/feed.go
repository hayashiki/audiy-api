package usecase

import (
	"context"
	"github.com/hayashiki/audiy-api/src/domain/repository"
	"log"

	"github.com/hayashiki/audiy-api/src/domain/model"
)

type FeedUsecase interface {
	GetConnection(ctx context.Context, userID string, cursor string, limit int, filter *model.FeedEvent, orderBy string) (*model.FeedConnection, error)
	Get(ctx context.Context, userID string, id string) (*model.Feed, error)
	Put(ctx context.Context, userID string, audioID string, event model.FeedEvent) (*model.Feed, error)
}

func NewFeedUsecase(feedRepo repository.FeedRepository, audioRepo repository.AudioRepository) FeedUsecase {
	return &feedUsecase{feedRepo: feedRepo, audioRepo: audioRepo}
}

type feedUsecase struct {
	feedRepo  repository.FeedRepository
	audioRepo repository.AudioRepository
}

func (u *feedUsecase) GetConnection(ctx context.Context, userID string, cursor string, limit int, filter *model.FeedEvent, orderBy string) (*model.FeedConnection, error) {
	var filters map[string]interface{}

	if filter == nil {
		filter = nil
	} else {
		switch *filter {
		case model.FeedEventLiked:
			filters = map[string]interface{}{
				"Liked": true,
			}
		case model.FeedEventUnliked:
			filters = map[string]interface{}{
				"Liked": false,
			}
		case model.FeedEventStared:
			filters = map[string]interface{}{
				"Stared": true,
			}
		case model.FeedEventUnstared:
			filters = map[string]interface{}{
				"Stared": false,
			}
		case model.FeedEventPlayed:
			filters = map[string]interface{}{
				"Played": true,
			}
		case model.FeedEventUnplayed:
			filters = map[string]interface{}{
				"Played": false,
			}
		case model.FeedEventAll:
			filters = map[string]interface{}{}
		default:
			filters = map[string]interface{}{}
		}
	}

	feeds, nextCursor, HasMore, err := u.feedRepo.GetAll(ctx, userID, filters, cursor, limit, orderBy)
	if err != nil {
		return nil, err
	}
	feedEdges := make([]*model.FeedEdge, len(feeds))
	for i, a := range feeds {
		feedEdges[i] = &model.FeedEdge{
			Cursor: nextCursor,
			Node:   a,
		}
	}
	return &model.FeedConnection{
		PageInfo: &model.PageInfo{
			Cursor:  nextCursor,
			HasMore: HasMore,
		},
		Edges: feedEdges,
	}, nil
}

func (u *feedUsecase) Get(ctx context.Context, userID string, id string) (*model.Feed, error) {
	return u.feedRepo.Get(ctx, userID, id)
}

func (u *feedUsecase) Put(ctx context.Context, userID string, audioID string, event model.FeedEvent) (*model.Feed, error) {
	feed, err := u.feedRepo.GetByAudio(ctx, userID, audioID)
	if err != nil {
		return nil, err
	}

	audio, err := u.audioRepo.Get(ctx, feed.AudioID)
	if err != nil {
		return nil, err
	}

	switch event {
	case model.FeedEventPlayed:
		audio.PlayCount += 1
		feed.Played = true
	case model.FeedEventUnplayed:
		if !feed.Played {
			return feed, nil
		}
		feed.Played = false
	case model.FeedEventStared:
		if feed.Stared {
			return feed, nil
		}
		feed.Stared = true
	case model.FeedEventUnstared:
		if !feed.Stared {
			return feed, nil
		}
		feed.Stared = false
	case model.FeedEventLiked:
		if feed.Liked {
			return feed, nil
		}
		feed.Liked = true
		audio.LikeCount += 1
	case model.FeedEventUnliked:
		if !feed.Liked {
			return feed, nil
		}
		feed.Liked = false
		audio.LikeCount -= 1
	default:
		log.Printf("invalid event %s", event)
	}
	if err := u.feedRepo.Put(ctx, userID, feed); err != nil {
		return nil, err
	}
	if err := u.audioRepo.Put(ctx, audio); err != nil {
		return nil, err
	}
	return feed, err
}
