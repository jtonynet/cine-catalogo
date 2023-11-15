package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"

	"github.com/jtonynet/cine-catalogo/handlers/requests"
	"github.com/jtonynet/cine-catalogo/handlers/responses"
	"github.com/jtonynet/cine-catalogo/internal/database"
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

				HATEOASProperties: responses.HATEOASProperties{
					Links: responses.HATEOASBaseLinks{
						Self: responses.HREFObject{
							HREF: fmt.Sprintf("http://localhost:8080/v1/addresses/%s", address.UUID.String()),
						},
					},
				},
			},
		)
	}

	addressesLinks := AddressesLinks{
		Self:            responses.HREFObject{HREF: "http://localhost:8080/v1/addresses"},
		CreateAddresses: responses.HREFObject{HREF: "http://localhost:8080/v1/addresses"},
	}

	addressResponseList := AddressResponseList{
		Addresses: &response,
	}

	resultEmbedded := responses.HATEOASResultEmbedded{
		Embedded: addressResponseList,
		Links:    addressesLinks,
	}

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"RetrieveAddressList",
		resultEmbedded,
		responses.HALHeaders,
	)
}

type AddressResponseList struct {
	Addresses *[]responses.Address `json:"addresses"`
}

type AddressesLinks struct {
	Self            responses.HREFObject `json:"self"`
	CreateAddresses responses.HREFObject `json:"create-addresses"`
}
