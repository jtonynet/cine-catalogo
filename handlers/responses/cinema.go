package responses

import "github.com/google/uuid"

type Cinema struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Capacity    int64     `json:"capacity"`
}
