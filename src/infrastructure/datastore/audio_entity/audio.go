package audio_entity

import (
	"github.com/hayashiki/audiy-api/src/domain/model"
	"time"
)

const kind = "Audio"

func onlyID(id string) *entity {
	return &entity{ID: id}
}

type entity struct {
	kind      string `boom:"kind,Audio"`
	ID        string `boom:"id"`
	Name      string
	Length    float64
	URL         string
	Mimetype     string
	Transcribed  bool
	LikeCount    int
	PlayCount    int
	CommentCount int
	PublishedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (f *entity) toDomain() *model.Audio {
	return &model.Audio{
		ID:        f.ID,
		Name:      f.Name,
		Length:     f.Length,
		//URL:     f.URL,
		Mimetype:     f.Mimetype,
		Transcribed:     f.Transcribed,
		LikeCount:     f.LikeCount,
		PlayCount:  f.PlayCount,
		CommentCount:  f.CommentCount,
		PublishedAt:  f.PublishedAt,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func toEntity(from *model.Audio) *entity {
	return &entity{
		ID:        from.ID,
		Name:      from.Name,
		Length:     from.Length,
		//URL:     f.URL,
		Mimetype:     from.Mimetype,
		Transcribed:     from.Transcribed,
		LikeCount:     from.LikeCount,
		PlayCount:  from.PlayCount,
		CommentCount:  from.CommentCount,
		PublishedAt:  from.PublishedAt,
		CreatedAt: from.CreatedAt,
		UpdatedAt: from.UpdatedAt,
	}
}
