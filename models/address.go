package models

import (
	"github.com/google/uuid"
)

type Address struct {
	BaseSchema
	UUID        uuid.UUID `gorm:"type:uuid;unique;not null"`
	Country     string
	State       string
	Telephone   string
	Description string
	PostalCode  string
	Name        string
}

func NewAddress(
	uuid uuid.UUID,
	Country string,
	State string,
	Telephone string,
	Description string,
	PostalCode string,
	Name string,
) (Address, error) {
	a := Address{
		UUID:        uuid,
		Country:     Country,
		State:       State,
		Telephone:   Telephone,
		Description: Description,
		PostalCode:  PostalCode,
		Name:        Name,
	}

	return a, nil
}
