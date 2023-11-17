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

	HATEOASListItemProperties
}

type HATEOASCinemaItemLinks struct {
	Self HREFObject `json:"self"`
}

type HATEOASCinemaListLinks struct {
	Self HREFObject `json:"self"`
}
type HATEOASCinemaList struct {
	Cinemas *[]Cinema `json:"cinemas"`
}

func NewCinema(c models.Cinema, baseURL string) Cinema {
	return Cinema{
		UUID:        c.UUID,
		Name:        c.Name,
		Description: c.Description,
		Capacity:    c.Capacity,

		HATEOASListItemProperties: HATEOASListItemProperties{
			Links: HATEOASCinemaItemLinks{
				Self: HREFObject{
					HREF: fmt.Sprintf("%s/cinemas/%s", baseURL, c.UUID.String()),
				},
			},
		},
	}
}
