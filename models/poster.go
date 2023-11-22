package models

import "github.com/google/uuid"

type Poster struct {
	BaseModel
	UUID            uuid.UUID `gorm:"type:uuid;unique;not null"`
	MovieID         uint
	Movie           Movie
	Name            string
	ContentType     string
	AlternativeText string
	Path            string
}

func NewPoster(
	UUID uuid.UUID,
	movieID uint,
	name,
	contentType,
	alternativeText,
	path string,
) (Poster, error) {
	p := Poster{
		UUID:            UUID,
		MovieID:         movieID,
		Name:            name,
		ContentType:     contentType,
		AlternativeText: alternativeText,
		Path:            path,
	}

	return p, nil
}
