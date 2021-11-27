package usecase

import (
	"context"
	"log"

	"github.com/hayashiki/audiy-api/src/domain/entity"
)

type FeedUsecase interface {
	GetConnection(ctx context.Context, userID string, cursor string, limit int, filter *entity.FeedEvent, order []string) (*entity.FeedConnection, error)
	Get(ctx context.Context, id int64, userID string) (*entity.Feed, error)
	Put(ctx context.Context, userID string, audioID string, event entity.FeedEvent) (*entity.Feed, error)
}

func NewFeedUsecase(feedRepo entity.FeedRepository, audioRepo entity.AudioRepository) FeedUsecase {
	return &feedUsecase{feedRepo: feedRepo, audioRepo: audioRepo}
}

type feedUsecase struct {
	feedRepo  entity.FeedRepository
	audioRepo entity.AudioRepository
}

func (u *feedUsecase) GetConnection(ctx context.Context, userID string, cursor string, limit int, filter *entity.FeedEvent, order []string) (*entity.FeedConnection, error) {
	var filters map[string]interface{}

	if filter == nil {
		filter = nil
	} else {
		switch *filter {
		case entity.FeedEventLiked:
			filters = map[string]interface{}{
				"liked": true,
			}
		case entity.FeedEventUnliked:
			filters = map[string]interface{}{
				"liked": false,
			}
		case entity.FeedEventStared:
			filters = map[string]interface{}{
				"stared": true,
			}
		case entity.FeedEventUnstared:
			filters = map[string]interface{}{
				"stared": false,
			}
		case entity.FeedEventPlayed:
			filters = map[string]interface{}{
				"played": true,
			}
		case entity.FeedEventUnplayed:
			filters = map[string]interface{}{
				"played": false,
			}
		case entity.FeedEventAll:
			filters = map[string]interface{}{}
		default:
			filters = map[string]interface{}{}
		}
	}

	feeds, nextCursor, HasMore, err := u.feedRepo.FindAll(ctx, userID, filters, cursor, limit, order...)
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
			HasMore: HasMore,
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

	audio, err := u.audioRepo.Find(ctx, feed.AudioKey.Name)
	if err != nil {
		return nil, err
	}

	switch event {
	case entity.FeedEventPlayed:
		audio.PlayCount += 1
		feed.Played = true
	case entity.FeedEventUnplayed:
		if !feed.Played {
			return feed, nil
		}
		feed.Played = false
	case entity.FeedEventStared:
		if feed.Stared {
			return feed, nil
		}
		feed.Stared = true
	case entity.FeedEventUnstared:
		if !feed.Stared {
			return feed, nil
		}
		feed.Stared = false
	case entity.FeedEventLiked:
		if feed.Liked {
			return feed, nil
		}
		feed.Liked = true
		audio.LikeCount += 1
	case entity.FeedEventUnliked:
		if !feed.Liked {
			return feed, nil
		}
		feed.Liked = false
		audio.LikeCount -= 1
	default:
		log.Printf("invalid event %s", event)
	}
	if err := u.feedRepo.Save(ctx, userID, feed); err != nil {
		return nil, err
	}
	if err := u.audioRepo.Save(ctx, audio); err != nil {
		return nil, err
	}
	return feed, err
}
