package transcript_entity

import (
	"time"

	"github.com/hayashiki/audiy-api/src/domain/model"
)

const kind = "Transcript"

func onlyID(id int64) *entity {
	return &entity{ID: id}
}

type entity struct {
	kind       string `boom:"kind,Transcript"`
	ID         int64  `boom:"id"`
	AudioID    string
	Body       string `datastore:",noindex"`
	Monologues []model.Monologue
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (f *entity) toDomain() *model.Transcript {
	return &model.Transcript{
		ID:         f.ID,
		//AudioID:    f.AudioID,
		Body:       f.Body,
		Monologues: f.Monologues,
		CreatedAt:  f.CreatedAt,
		UpdatedAt:  f.UpdatedAt,
	}
}

func toEntity(from *model.Transcript) *entity {
	return &entity{
		ID:         from.ID,
		AudioID:    from.AudioID,
		Body:       from.Body,
		Monologues: from.Monologues,
		CreatedAt:  from.CreatedAt,
		UpdatedAt:  from.UpdatedAt,
	}
}

//
//type Monologue struct {
//	Elements []MonologueElement `json:"elements" datastore:"elements"`
//}
//
//type MonologueElement struct {
//	StartTime  float64 `json:"start_time" datastore:"start_time"`
//	EndTime    float64 `json:"end_time" datastore:"end_time"`
//	Word       string `json:"word" datastore:"word"`
//	WordKana   string `json:"word_kana" datastore:"word_kana"`
//	Confidence float32 `json:"confidence" datastore:"confidence"`
//}
//
//type Transcript struct {
//	Key        *datastore.Key `datastore:"__key__"`
//	ID         int64          `json:"id" datastore:"-"`
//	AudioKey   *datastore.Key `json:"audio_key" datastore:"audio_key"`
//	Body       string         `json:"body" datastore:"body,noindex"`
//	Monologues []Monologue `json:"monologues" datastore:"monologues"`
//	CreatedAt  time.Time `json:"createdAt" datastore:"created_at"`
//	UpdatedAt  time.Time `json:"updatedAt" datastore:"updated_at"`
//}
