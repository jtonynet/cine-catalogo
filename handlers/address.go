package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"

	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/handlers/requests"
	"github.com/jtonynet/cine-catalogo/handlers/responses"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/internal/hateoas"
	"github.com/jtonynet/cine-catalogo/models"
)

func CreateAddresses(ctx *gin.Context) {
	var requestList []requests.Address
	if err := ctx.ShouldBindBodyWith(&requestList, binding.JSON); err != nil {

		var singleRequest requests.Address
		if err := ctx.ShouldBindBodyWith(&singleRequest, binding.JSON); err != nil {
			// TODO: Implements in future
			return
		}

		requestList = append(requestList, singleRequest)
	}

	var addressList []models.Address
	for _, request := range requestList {
		address, err := models.NewAddress(
			uuid.New(),
			request.Country,
			request.State,
			request.Telephone,
			request.Description,
			request.PostalCode,
			request.Name,
		)
		if err != nil {
			// TODO Implements
			return
		}

		addressList = append(addressList, address)
	}

	if err := database.DB.Create(&addressList).Error; err != nil {
		responses.SendError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	responseList := []responses.Address{}
	for _, address := range addressList {
		responseList = append(responseList,
			responses.Address{
				UUID:        address.UUID,
				Country:     address.Country,
				State:       address.State,
				Telephone:   address.Telephone,
				Description: address.Description,
				PostalCode:  address.PostalCode,
				Name:        address.Name,
			},
		)
	}

	responses.SendSuccess(ctx, http.StatusOK, "CreateAddresses", responseList, responses.HALHeaders)
}

func RetrieveAddress(ctx *gin.Context) {
	uuid := uuid.MustParse(ctx.Param("addressId"))

	address := models.Address{UUID: uuid}
	database.DB.Where(&models.Address{UUID: uuid}).First(&address)

	response := responses.Address{
		UUID:        address.UUID,
		Country:     address.Country,
		State:       address.State,
		Telephone:   address.Telephone,
		Description: address.Description,
		PostalCode:  address.PostalCode,
		Name:        address.Name,
	}

	responses.SendSuccess(ctx, http.StatusOK, "RetrieveAddress", response, nil)
}

func RetrieveAddressList(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	rootURL := cfg.Host

	addresses := []models.Address{}

	if err := database.DB.Find(&addresses).Error; err != nil {
		// TODO: Implements in future
		return
	}

	response := []responses.Address{}
	for _, address := range addresses {

		response = append(
			response,
			responses.Address{
				UUID:        address.UUID,
				Country:     address.Country,
				State:       address.State,
				Telephone:   address.Telephone,
				Description: address.Description,
				PostalCode:  address.PostalCode,
				Name:        address.Name,

				HATEOASListItemProperties: responses.HATEOASListItemProperties{
					Links: HATEOASAddressesLinks{
						Self: responses.HREFObject{
							HREF: fmt.Sprintf("%s/addresses/%s", rootURL, address.UUID.String()),
						},
						CreateAddressesCinemas: responses.HREFObject{
							HREF: fmt.Sprintf("%s/addresses/%s/cinemas", rootURL, address.UUID.String()),
						},
					},
				},
			},
		)
	}

	addressesLinks := AddressesLinks{
		Self:            responses.HREFObject{HREF: fmt.Sprintf("%s/addresses", rootURL)},
		CreateAddresses: responses.HREFObject{HREF: fmt.Sprintf("%s/addresses", rootURL)},
	}

	addressResponseList := AddressResponseList{
		Addresses: &response,
	}

	root := hateoas.NewRoot(rootURL)
	createAddressesPost, err := hateoas.NewResource(
		"create-addresses",
		fmt.Sprintf("%s/%s", rootURL, "addresses"),
		http.MethodPost,
	)
	if err != nil {
		// TODO: implements on future
		return
	}
	createAddressesPost.RequestToProperties(requests.Address{})
	root.AddResource(createAddressesPost)

	createCinemasUrl := fmt.Sprintf("%s/addresses/:addressId/cinemas", rootURL)
	createCinemasPost, err := hateoas.NewResource(
		"create-addresses-cinemas",
		createCinemasUrl,
		http.MethodPost,
	)
	if err != nil {
		// TODO: implements on future
		return
	}
	createCinemasPost.RequestToProperties(requests.Cinema{})
	root.AddResource(createCinemasPost)

	rootEncoded, err := root.Encode()
	if err != nil {
		// TODO: implements on future
		return
	}

	templateString := gjson.Get(string(rootEncoded), "_templates").String()
	var templateJSON interface{}
	json.Unmarshal([]byte(templateString), &templateJSON)

	resultEmbedded := responses.HATEOASResultEmbedded{
		Embedded:  addressResponseList,
		Links:     addressesLinks,
		Templates: templateJSON,
	}

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"retrieve-address-list",
		resultEmbedded,
		responses.HALHeaders,
	)
}

type HATEOASAddressesLinks struct {
	Self                   responses.HREFObject `json:"self"`
	CreateAddressesCinemas responses.HREFObject `json:"create-addresses-cinemas"`
}

type AddressResponseList struct {
	Addresses *[]responses.Address `json:"addresses"`
}

type AddressesLinks struct {
	Self            responses.HREFObject `json:"self"`
	CreateAddresses responses.HREFObject `json:"create-addresses"`
}
