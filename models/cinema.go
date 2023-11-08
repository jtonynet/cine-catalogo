package models

import "github.com/google/uuid"

type Cinema struct {
	BaseModel
	UUID        uuid.UUID `gorm:"type:uuid;unique;not null"`
	Name        string
	Description string
	Capacity    int64
}

func NewCinema(
	UUID uuid.UUID,
	Name string,
	Description string,
	Capacity int64,
) (Cinema, error) {
	c := Cinema{
		UUID:        UUID,
		Name:        Name,
		Description: Description,
		Capacity:    Capacity,
	}

	return c, nil
}
