package user_entity

import (
	"go.mercari.io/datastore"
	"time"

	"github.com/hayashiki/audiy-api/src/domain/model"
)

const kind = "User"

func onlyID(id string) *entity {
	return &entity{ID: id}
}

type entity struct {
	kind      string `boom:"kind,User"`
	ID        string `boom:"id"`
	Key       datastore.Key `datastore:"__key__"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	PhotoURL  string `json:"photoURL"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (f *entity) toDomain() *model.User {
	return &model.User{
		ID:        f.ID,
		Key2:       f.Key,
		Email:     f.Email,
		Name:      f.Name,
		PhotoURL:  f.PhotoURL,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func toEntity(from *model.User) *entity {
	return &entity{
		ID:        from.ID,
		Email:     from.Email,
		Key:       from.Key2,
		Name:      from.Name,
		PhotoURL:  from.PhotoURL,
		CreatedAt: from.CreatedAt,
		UpdatedAt: from.UpdatedAt,
	}
}
