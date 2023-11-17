package responses

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/models"
)

type Address struct {
	UUID        uuid.UUID `json:"uuid"`
	Country     string    `json:"country"`
	State       string    `json:"state"`
	Telephone   string    `json:"telephone"`
	Description string    `json:"description"`
	PostalCode  string    `json:"postalCode"`
	Name        string    `json:"name"`

	HATEOASListItemProperties
}

type HATEOASAddressItemLinks struct {
	Self                   HREFObject `json:"self"`
	CreateAddressesCinemas HREFObject `json:"create-addresses-cinemas"`
	RetrieveCinemaList     HREFObject `json:"retrieve-cinema-list"`
}

type HATEOASAddressListLinks struct {
	Self            HREFObject `json:"self"`
	CreateAddresses HREFObject `json:"create-addresses"`
}
type HATEOASAddressList struct {
	Addresses *[]Address `json:"addresses"`
}

func NewAddress(a models.Address, baseURL string) Address {
	return Address{
		UUID:        a.UUID,
		Country:     a.Country,
		State:       a.State,
		Telephone:   a.Telephone,
		Description: a.Description,
		PostalCode:  a.PostalCode,
		Name:        a.Name,

		HATEOASListItemProperties: HATEOASListItemProperties{
			Links: HATEOASAddressItemLinks{
				Self: HREFObject{
					HREF: fmt.Sprintf("%s/addresses/%s", baseURL, a.UUID.String()),
				},
				CreateAddressesCinemas: HREFObject{
					HREF: fmt.Sprintf("%s/addresses/%s/cinemas", baseURL, a.UUID.String()),
				},
				RetrieveCinemaList: HREFObject{
					HREF: fmt.Sprintf("%s/addresses/%s/cinemas", baseURL, a.UUID.String()),
				},
			},
		},
	}
}
