package feed_entity

import (
	"cloud.google.com/go/datastore"
	"github.com/hayashiki/audiy-api/src/domain/model"
	"time"
)

const kind = "Feed"
const parentKind = "User"

func onlyID(id string) *entity {
	return &entity{ID: id}
}

type parent struct {
	kind      string `boom:"kind,Audio"`
	ID        string `boom:"id"`
}

type entity struct {
	Key         *datastore.Key `datastore:"__key__"`
	//kind      string `boom:"kind,Feed"`
	ID        string `datastore:"-"` // boom:"id"
	AudioID  string `datastore:"AudioID"`
	UserID   string `datastore:"UserID"`
	Played  bool `datastore:"Played"`
	Liked  bool `datastore:"Liked"`
	Stared  bool `datastore:"Stared"`
	StartTime float64 `datastore:"StartTime"`
	PublishedAt time.Time `datastore:"PublishedAt"`
	CreatedAt time.Time `datastore:"CreatedAt"`
	UpdatedAt time.Time `datastore:"UpdatedAt"`
}

func (f *entity) toDomain() *model.Feed {
	return &model.Feed{
		//ID:        f.ID,
		AudioID:      f.AudioID,
		UserID:       f.UserID,
		Played:     f.Played,
		Liked:     f.Liked,
		Stared:     f.Stared,
		StartTime:     f.StartTime,
		PublishedAt:  f.PublishedAt,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func toEntity(from *model.Feed) *entity {
	return &entity{
		ID:        string(from.ID()),
		//ParentKey:  from.ParentKey,
		AudioID:      from.AudioID,
		UserID:       from.UserID,
		Played:     from.Played,
		Liked:     from.Liked,
		Stared:     from.Stared,
		StartTime:     0,
		PublishedAt:  from.PublishedAt,
		CreatedAt: from.CreatedAt,
		UpdatedAt: from.UpdatedAt,
	}
}
