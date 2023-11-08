package models

import "github.com/google/uuid"

type Movie struct {
	BaseModel
	UUID        uuid.UUID `gorm:"type:uuid;unique;not null"`
	Name        string
	Description string
	AgeRating   int64
	Subtitled   bool
	Poster      string
}

func NewMovie(
	UUID uuid.UUID,
	Name string,
	Description string,
	AgeRating int64,
	Subtitled bool,
	Poster string,
) (Movie, error) {
	m := Movie{
		UUID:        UUID,
		Name:        Name,
		Description: Description,
		AgeRating:   AgeRating,
		Subtitled:   Subtitled,
		Poster:      Poster,
	}

	return m, nil
}
