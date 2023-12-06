package responses

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/models"
)

type baseCinema struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Capacity    int64     `json:"capacity"`
}

type Cinema struct {
	baseCinema

	Links HATEOASCinemasLinks `json:"_links,omitempty"`

	Templates interface{} `json:"_templates,omitempty"`
}

type HATEOASCinemasLinks struct {
	Self         HATEOASLink `json:"self"`
	UpdateCinema HATEOASLink `json:"update-cinema"`
	DeleteCinema HATEOASLink `json:"delete-cinema"`
	Address      HATEOASLink `json:"address"`
}

func NewCinema(
	model models.Cinema,
	addressLink,
	baseURL string,
	templates interface{},
) Cinema {
	cinema := Cinema{
		baseCinema: baseCinema{
			UUID:        model.UUID,
			Name:        model.Name,
			Description: model.Description,
			Capacity:    model.Capacity,
		},

		Links: HATEOASCinemasLinks{
			Self:         HATEOASLink{HREF: fmt.Sprintf("%s/cinemas/%s", baseURL, model.UUID.String())},
			UpdateCinema: HATEOASLink{HREF: fmt.Sprintf("%s/cinemas/%s", baseURL, model.UUID.String())},
			DeleteCinema: HATEOASLink{HREF: fmt.Sprintf("%s/cinemas/%s", baseURL, model.UUID.String())},
			Address:      HATEOASLink{HREF: addressLink},
		},

		Templates: templates,
	}

	return cinema
}

type HATEOASCinemaListLinks struct {
	Self HATEOASLink `json:"self"`
}
type HATEOASCinemaList struct {
	Cinemas []Cinema `json:"cinemas,omitempty"`
}
