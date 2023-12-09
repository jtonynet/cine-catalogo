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

type AddressOption func(*Address)

func NewAddress(
	model models.Address,
	baseURL string,
	options ...AddressOption,
) Address {
	address := Address{

		UUID:        model.UUID,
		Country:     model.Country,
		State:       model.State,
		Telephone:   model.Telephone,
		Description: model.Description,
		PostalCode:  model.PostalCode,
		Name:        model.Name,

		Links: HATEOASAddressLinks{
			Self:                   HATEOASLink{HREF: fmt.Sprintf("%s/addresses/%s", baseURL, model.UUID.String())},
			CreateAddressesCinemas: HATEOASLink{HREF: fmt.Sprintf("%s/addresses/%s/cinemas", baseURL, model.UUID.String())},
			RetrieveCinemaList:     HATEOASLink{HREF: fmt.Sprintf("%s/addresses/%s/cinemas", baseURL, model.UUID.String())},
			UpdateAddress:          HATEOASLink{HREF: fmt.Sprintf("%s/addresses/%s", baseURL, model.UUID.String())},
			DeleteAddress:          HATEOASLink{HREF: fmt.Sprintf("%s/addresses/%s", baseURL, model.UUID.String())},
		},
	}

	for _, opt := range options {
		opt(&address)
	}

	return address
}

func WithAddressTemplates(templates interface{}) AddressOption {
	return func(a *Address) {
		a.Templates = templates
	}
}

type HATEOASAddressListLinks struct {
	Self            HATEOASLink `json:"self"`
	CreateAddresses HATEOASLink `json:"create-addresses"`
}
type HATEOASAddressList struct {
	Addresses []Address `json:"addresses"`
}
