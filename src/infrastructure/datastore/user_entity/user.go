package user_entity

import (
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
	Email     string
	Name      string
	PhotoURL  string
	ProviderID  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e *entity) toDomain() *model.User {
	return &model.User{
		ID:       e.ID,
		Email:    e.Email,
		Name:     e.Name,
		PhotoURL: e.PhotoURL,
		ProviderID: e.ProviderID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func toEntity(from *model.User) *entity {
	return &entity{
		ID:        from.ID,
		Email:     from.Email,
		Name:      from.Name,
		PhotoURL:  from.PhotoURL,
		ProviderID: from.ProviderID,
		CreatedAt: from.CreatedAt,
		UpdatedAt: from.UpdatedAt,
	}
}
