package entity

import "time"

type Like struct {
	ID        string    `json:"id"`
	User      *User     `json:"user"`
	Audio     *Audio    `json:"audio"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Like) IsNode() {}