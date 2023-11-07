package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Theater struct {
	gorm.Model
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	Capacity int64     `json:"capacity"`
}
