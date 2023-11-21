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

	HATEOASListItemResult
}

type HATEOASAddressItemLinks struct {
	Self                   HATEOASLink `json:"self"`
	CreateAddressesCinemas HATEOASLink `json:"create-addresses-cinemas"`
	RetrieveCinemaList     HATEOASLink `json:"retrieve-cinema-list"`
}

type HATEOASAddressListLinks struct {
	Self            HATEOASLink `json:"self"`
	CreateAddresses HATEOASLink `json:"create-addresses"`
}
type HATEOASAddressList struct {
	Addresses *[]Address `json:"addresses"`
}

type AddressOption func(*Address)

func NewAddress(
	model models.Address,
	baseURL string,
	options ...AddressOption,
) *Address {
	address := &Address{
		UUID:        model.UUID,
		Country:     model.Country,
		State:       model.State,
		Telephone:   model.Telephone,
		Description: model.Description,
		PostalCode:  model.PostalCode,
		Name:        model.Name,

		HATEOASListItemResult: HATEOASListItemResult{
			Links: HATEOASAddressItemLinks{
				Self: HATEOASLink{
					HREF: fmt.Sprintf("%s/addresses/%s", baseURL, model.UUID.String()),
				},
				CreateAddressesCinemas: HATEOASLink{
					HREF: fmt.Sprintf("%s/addresses/%s/cinemas", baseURL, model.UUID.String()),
				},
				RetrieveCinemaList: HATEOASLink{
					HREF: fmt.Sprintf("%s/addresses/%s/cinemas", baseURL, model.UUID.String()),
				},
			},
		},
	}

	for _, opt := range options {
		opt(address)
	}

	return address
}

func WithAddressTemplates(templates interface{}) AddressOption {
	return func(address *Address) {
		address.Templates = templates
	}
}
