package feed_entity

import (
	"github.com/hayashiki/audiy-api/src/domain/model"
	"go.mercari.io/datastore"
	"time"
)

const kind = "Feed"
const parentKind = "User"

func onlyID(id int64) *entity {
	return &entity{ID: id}
}

type parent struct {
	kind      string `boom:"kind,Audio"`
	ID        string `boom:"id"`
}

type entity struct {
	kind      string `boom:"kind,Feed"`
	ParentKey datastore.Key `datastore:"-" boom:"parent"`
	// parent(userID) + audioに変更する
	ID        int64 `boom:"id"`
	AudioID  string
	Played  bool
	Liked  bool
	Stared  bool
	StartTime float64
	PublishedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (f *entity) toDomain() *model.Feed {
	return &model.Feed{
		ID:        f.ID,
		//ParentKey:  f.ParentKey,
		AudioID:      f.AudioID,
		Played:     f.Played,
		Liked:     f.Liked,
		Stared:     f.Stared,
		//StartTime:     &f.StartTime,
		PublishedAt:  f.PublishedAt,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func toEntity(from *model.Feed) *entity {

	if from.StartTime == nil {
		//*from.StartTime = 0
	}

	return &entity{
		ID:        from.ID,
		//ParentKey:  from.ParentKey,
		AudioID:      from.AudioID,
		Played:     from.Played,
		Liked:     from.Liked,
		Stared:     from.Stared,
		//StartTime:     *from.StartTime,
		PublishedAt:  from.PublishedAt,
		CreatedAt: from.CreatedAt,
		UpdatedAt: from.UpdatedAt,
	}
}
