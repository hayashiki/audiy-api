package comment_entity

import (
	"github.com/hayashiki/audiy-api/src/domain/model"
	"time"
)

const kind = "Comment"

func onlyID(id int64) *entity {
	return &entity{ID: id}
}

type entity struct {
	kind      string `boom:"kind,Comment"`
	ID        int64 `boom:"id"`
	UserID    string
	AudioID   string
	Body      string `datastore:",noindex"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (f *entity) toDomain() *model.Comment {
	return &model.Comment{
		ID:        f.ID,
		UserID:    f.UserID,
		AudioID:   f.AudioID,
		Body:      f.Body,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func toEntity(from *model.Comment) *entity {
	return &entity{
		ID:        from.ID,
		UserID:    from.UserID,
		AudioID:   from.AudioID,
		Body:      from.Body,
		CreatedAt: from.CreatedAt,
		UpdatedAt: from.UpdatedAt,
	}
}

