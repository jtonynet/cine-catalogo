package responses

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/models"
)

type baseAddress struct {
	UUID        uuid.UUID `json:"uuid"`
	Country     string    `json:"country"`
	State       string    `json:"state"`
	Telephone   string    `json:"telephone"`
	Description string    `json:"description"`
	PostalCode  string    `json:"postalCode"`
	Name        string    `json:"name"`
}

type Address struct {
	baseAddress

	Links HATEOASAddressLinks `json:"_links,omitempty"`

	Templates interface{} `json:"_templates,omitempty"`
}

type HATEOASAddressLinks struct {
	Self                   HATEOASLink `json:"self"`
	CreateAddressesCinemas HATEOASLink `json:"create-addresses-cinemas"`
	RetrieveCinemaList     HATEOASLink `json:"retrieve-cinema-list"`
	UpdateAddress          HATEOASLink `json:"update-address"`
	DeleteAddress          HATEOASLink `json:"delete-address"`
}

func NewAddress(
	model models.Address,
	baseURL string,
	templates interface{},
) Address {
	address := Address{
		baseAddress: baseAddress{
			UUID:        model.UUID,
			Country:     model.Country,
			State:       model.State,
			Telephone:   model.Telephone,
			Description: model.Description,
			PostalCode:  model.PostalCode,
			Name:        model.Name,
		},

		Links: HATEOASAddressLinks{
			Self:                   HATEOASLink{HREF: fmt.Sprintf("%s/addresses/%s", baseURL, model.UUID.String())},
			CreateAddressesCinemas: HATEOASLink{HREF: fmt.Sprintf("%s/addresses/%s/cinemas", baseURL, model.UUID.String())},
			RetrieveCinemaList:     HATEOASLink{HREF: fmt.Sprintf("%s/addresses/%s/cinemas", baseURL, model.UUID.String())},
			UpdateAddress:          HATEOASLink{HREF: fmt.Sprintf("%s/addresses/%s", baseURL, model.UUID.String())},
			DeleteAddress:          HATEOASLink{HREF: fmt.Sprintf("%s/addresses/%s", baseURL, model.UUID.String())},
		},

		Templates: templates,
	}

	return address
}

type HATEOASAddressListLinks struct {
	Self            HATEOASLink `json:"self"`
	CreateAddresses HATEOASLink `json:"create-addresses"`
}
type HATEOASAddressList struct {
	Addresses []Address `json:"addresses"`
}
