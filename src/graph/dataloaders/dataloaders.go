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


//func Middleware(repo repository.Repository, next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
//			ProductUsersByIDs: UserLoader{
//				maxBatch: 100,
//				wait:     15 * time.Millisecond,
//				fetch: func(ids []int) ([]*model.User, []error) {
//					data, err := repo.GetLoadUsers(r.Context(), ids)
//					if err != nil {
//						return nil, []error{err}
//					}
//					if len(data) == 0 {
//						return nil, nil
//					}
//
//					list := make([]*model.User, len(data))
//					for i, d := range data {
//						user := &model.User{}
//						user.User = d
//						list[i] = user
//					}
//					return list, nil
//				},
//			},
//			ProductCategoryByIDs: CategoryLoader{
//				maxBatch: 100,
//				wait:     15 * time.Millisecond,
//				fetch: func(ids []int) ([]*model.Category, []error) {
//					data, err := repo.GetLoadCategories(r.Context(), ids)
//					if err != nil {
//						return nil, []error{err}
//					}
//					if len(data) == 0 {
//						return nil, nil
//					}
//
//					list := make([]*model.Category, len(data))
//					for i, d := range data {
//						category := &model.Category{}
//						category.Category = d
//						list[i] = category
//					}
//					return list, nil
//				},
//			},
//			ProductGenresByIDs: GenreLoader{
//				maxBatch: 100,
//				wait:     15 * time.Millisecond,
//				fetch: func(ids []int) ([][]*model.Genre, []error) {
//					dataList, err := repo.GetProductGenresByProductIDs(r.Context(), ids)
//					if err != nil {
//						return nil, []error{}
//					}
//					if len(dataList) == 0 {
//						return nil, nil
//					}
//
//					result := make([][]*model.Genre, len(ids))
//					aMap := make(map[int][]*model.Genre, len(ids))
//					for _, data := range dataList {
//						aMap[data.ProductID] = append(aMap[data.ProductID], &model.Genre{Genre: data.R.Genre})
//					}
//					for i, id := range ids {
//						result[i] = aMap[id]
//					}
//					return result, nil
//				},
//			},



