package models

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type IDRequest struct {
	ID uint `form:"id" json:"id" uri:"id"`
}

type DeleteRequest struct {
	IDList []uint `json:"idList"`
}
