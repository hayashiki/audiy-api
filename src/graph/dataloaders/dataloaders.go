package dataloaders

import (
	"context"
	"errors"
	"github.com/hayashiki/audiy-api/src/domain/entity"
	"net/http"
	"time"
)

//go:generate dataloaden AudioLoader string *github.com/hayashiki/audiy-api/src/domain/entity.Audio
type contextKey string

var loadersKey = contextKey("dataLoaders")

type loaders struct {
	AudioByID    *AudioLoader
}

func (s *DataLoaderService) Initialize(ctx context.Context) context.Context {
	return context.WithValue(ctx, loadersKey, &loaders{
		AudioByID:   newAudioByID(ctx, s.audioRepo),
	})
}

type DataLoaderService struct {
	audioRepo entity.AudioRepository
}

func NewDataLoaderService(audioRepo entity.AudioRepository) DataLoaderService {
	return DataLoaderService{audioRepo: audioRepo}
}

func (s *DataLoaderService) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		nextCtx := s.Initialize(ctx)
		r = r.WithContext(nextCtx)
		next.ServeHTTP(w, r)
	})
}

func (s *DataLoaderService) retrieve(ctx context.Context) (*loaders, error) {
	l, ok := ctx.Value(loadersKey).(*loaders)
	if !ok {
		return nil, errors.New("invalid")
	}
	return l, nil
}

func (s *DataLoaderService) AudioGetByID(ctx context.Context, id string) (*entity.Audio, error) {
	l, err := s.retrieve(ctx)
	if err != nil {
		return nil, err
	}
	return l.AudioByID.Load(id)
}

func newAudioByID(ctx context.Context, repo entity.AudioRepository) *AudioLoader {
	return NewAudioLoader(AudioLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(ids []string) ([]*entity.Audio, []error) {
			res, err := repo.GetMulti(ctx, ids)
			if err != nil {
				return nil, []error{err}
			}
			groupByID := make(map[string]*entity.Audio, len(ids))
			for _, r := range res {
				groupByID[r.Key.Name] = r
			}
			result := make([]*entity.Audio, len(ids))
			for i, id := range ids {
				result[i] = groupByID[id]
			}
			return result, nil
		},
	})
}
