package responses

import "github.com/google/uuid"

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
}

type HATEOASAddressListLinks struct {
	Self            HREFObject `json:"self"`
	CreateAddresses HREFObject `json:"create-addresses"`
}
type HATEOASAddressList struct {
	Addresses *[]Address `json:"addresses"`
}
