package fcm_entity

import "fmt"

const kind = "fcms"

func onlyID(id string) *fcm {
	return &fcm{ID: id}
}

type fcm struct {
	kind    string `boom:"kind,fcms"`
	ID       string `boom:"id"`
	UserID   string
	DeviceID string
	Token    string
}

func NewEntity(userID string, deviceID string, token string) *fcm {
	return &fcm{
		ID:       fmt.Sprintf("%s_%s", userID, deviceID),
		UserID:   userID,
		DeviceID: deviceID,
		Token:    token,
	}
}

