package models

import "github.com/google/uuid"

type Movie struct {
	BaseModel

	UUID        uuid.UUID `gorm:"type:uuid;unique;not null"`
	Name        string
	Description string
	AgeRating   int64
	Published   bool
	Subtitled   bool
	Posters     []Poster
}

func NewMovie(
	UUID uuid.UUID,
	Name string,
	Description string,
	AgeRating int64,
	Published,
	Subtitled bool,
) (Movie, error) {
	m := Movie{
		UUID:        UUID,
		Name:        Name,
		Description: Description,
		AgeRating:   AgeRating,
		Published:   Published,
		Subtitled:   Subtitled,
	}

	return m, nil
}
