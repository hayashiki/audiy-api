package fcm_entity

import (
	"fmt"
	"time"

	"github.com/hayashiki/audiy-api/src/domain/model"
)

const kind = "Fcm"

func onlyID(id string) *entity {
	return &entity{ID: id}
}

type entity struct {
	_kind      string `boom:"kind,Fcm"`
	ID        string `boom:"id" datastore:"-"`
	UserID    string
	DeviceID  string
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEntity(userID string, deviceID string, token string) *entity {
	return &entity{
		ID:       fmt.Sprintf("%s_%s", userID, deviceID),
		UserID:   userID,
		DeviceID: deviceID,
		Token:    token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (f *entity) toDomain() *model.Fcm {
	return &model.Fcm{
		ID:        f.ID,
		UserID:    f.UserID,
		DeviceID:  f.DeviceID,
		Token:     f.Token,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func toEntity(from *model.Fcm) *entity {
	return &entity{
		ID:        from.ID,
		UserID:    from.UserID,
		DeviceID:  from.DeviceID,
		Token:     from.Token,
		CreatedAt: from.CreatedAt,
		UpdatedAt: from.UpdatedAt,
	}
}
