package responses

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/models"
)

type Cinema struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Capacity    int64     `json:"capacity"`

	HATEOASListItemResult
}

type HATEOASCinemasItemLinks struct {
	Self HATEOASLink `json:"self"`
}

type HATEOASCinemaListLinks struct {
	Self HATEOASLink `json:"self"`
}
type HATEOASCinemaList struct {
	Cinemas *[]Cinema `json:"cinemas"`
}

type CinemaOption func(*Cinema)

func NewCinema(
	model models.Cinema,
	baseURL string,
	options ...CinemaOption,
) *Cinema {
	cinema := &Cinema{
		UUID:        model.UUID,
		Name:        model.Name,
		Description: model.Description,
		Capacity:    model.Capacity,

		HATEOASListItemResult: HATEOASListItemResult{
			Links: HATEOASCinemasItemLinks{
				Self: HATEOASLink{
					HREF: fmt.Sprintf("%s/cinemas/%s", baseURL, model.UUID.String()),
				},
			},
		},
	}

	for _, opt := range options {
		opt(cinema)
	}

	return cinema
}

func WithCinemaTemplates(templates interface{}) CinemaOption {
	return func(cinema *Cinema) {
		cinema.Templates = templates
	}
}
