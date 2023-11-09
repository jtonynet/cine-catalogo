package models

import "github.com/google/uuid"

/*
https://gorm.io/docs/belongs_to.html#Override-Foreign-Key
*/
type Cinema struct {
	BaseModel
	UUID        uuid.UUID `gorm:"type:uuid;unique;not null"`
	AddressID   uint
	Address     Address
	Name        string
	Description string
	Capacity    int64
}

func NewCinema(
	UUID uuid.UUID,
	AddressID uint,
	Name string,
	Description string,
	Capacity int64,
) (Cinema, error) {
	c := Cinema{
		UUID:        UUID,
		AddressID:   AddressID,
		Name:        Name,
		Description: Description,
		Capacity:    Capacity,
	}

	return c, nil
}
