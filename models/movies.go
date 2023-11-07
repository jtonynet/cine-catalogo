package models

import "github.com/google/uuid"

type Movies struct {
	BaseSchema
	UUID        uuid.UUID `gorm:"type:uuid;unique"`
	Description string
	Age_rating  string
	Subtitled   string
	Poster      string
}
