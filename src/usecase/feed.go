package usecase

import (
	"context"
	"log"

	entity2 "github.com/hayashiki/audiy-api/src/domain/entity"
)

type FeedUsecase interface {
	GetConnection(ctx context.Context, userID string, cursor string, limit int, filter *entity2.FeedEvent, order []string) (*entity2.FeedConnection, error)
	Get(ctx context.Context, id int64, userID string) (*entity2.Feed, error)
	Put(ctx context.Context, userID string, audioID string, event entity2.FeedEvent) (*entity2.Feed, error)
}

func NewFeedUsecase(feedRepo entity2.FeedRepository, audioRepo entity2.AudioRepository) FeedUsecase {
	return &feedUsecase{feedRepo: feedRepo, audioRepo: audioRepo}
}

type feedUsecase struct {
	feedRepo  entity2.FeedRepository
	audioRepo entity2.AudioRepository
}

func (u *feedUsecase) GetConnection(ctx context.Context, userID string, cursor string, limit int, filter *entity2.FeedEvent, order []string) (*entity2.FeedConnection, error) {
	var filters map[string]interface{}

	if filter == nil {
		filter = nil
	} else {
		switch *filter {
		case entity2.FeedEventLiked:
			filters = map[string]interface{}{
				"liked": true,
			}
		case entity2.FeedEventUnliked:
			filters = map[string]interface{}{
				"liked": false,
			}
		case entity2.FeedEventStared:
			filters = map[string]interface{}{
				"stared": true,
			}
		case entity2.FeedEventUnstared:
			filters = map[string]interface{}{
				"stared": false,
			}
		case entity2.FeedEventPlayed:
			filters = map[string]interface{}{
				"played": true,
			}
		case entity2.FeedEventUnplayed:
			filters = map[string]interface{}{
				"played": false,
			}
		case entity2.FeedEventAll:
			filters = map[string]interface{}{}
		default:
			filters = map[string]interface{}{}
		}
	}

	feeds, nextCursor, HasMore, err := u.feedRepo.FindAll(ctx, userID, filters, cursor, limit, order...)
	if err != nil {
		return nil, err
	}
	feedEdges := make([]*entity2.FeedEdge, len(feeds))
	for i, a := range feeds {
		feedEdges[i] = &entity2.FeedEdge{
			Cursor: nextCursor,
			Node:   a,
		}
	}
	return &entity2.FeedConnection{
		PageInfo: &entity2.PageInfo{
			Cursor:  nextCursor,
			HasMore: HasMore,
		},
		Edges: feedEdges,
	}, nil
}

func (u *feedUsecase) Get(ctx context.Context, id int64, userID string) (*entity2.Feed, error) {
	return u.feedRepo.Find(ctx, id, userID)
}

func (u *feedUsecase) Put(ctx context.Context, userID string, audioID string, event entity2.FeedEvent) (*entity2.Feed, error) {
	feed, err := u.feedRepo.FindByAudio(ctx, userID, audioID)
	if err != nil {
		return nil, err
	}

	audio, err := u.audioRepo.Find(ctx, feed.AudioKey.Name)
	if err != nil {
		return nil, err
	}

	switch event {
	case entity2.FeedEventPlayed:
		audio.PlayCount += 1
		feed.Played = true
	case entity2.FeedEventUnplayed:
		if !feed.Played {
			return feed, nil
		}
		feed.Played = false
	case entity2.FeedEventStared:
		if feed.Stared {
			return feed, nil
		}
		feed.Stared = true
	case entity2.FeedEventUnstared:
		if !feed.Stared {
			return feed, nil
		}
		feed.Stared = false
	case entity2.FeedEventLiked:
		if feed.Liked {
			return feed, nil
		}
		feed.Liked = true
		audio.LikeCount += 1
	case entity2.FeedEventUnliked:
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
