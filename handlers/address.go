package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"

	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/handlers/requests"
	"github.com/jtonynet/cine-catalogo/handlers/responses"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/models"
)

func CreateAddresses(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

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
			responses.NewAddress(address, versionURL),
		)
	}

	responses.SendSuccess(ctx, http.StatusOK, "create-addresses", responseList, responses.HALHeaders)
}

func RetrieveAddress(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	uuid := uuid.MustParse(ctx.Param("addressId"))

	address := models.Address{UUID: uuid}
	database.DB.Where(&models.Address{UUID: uuid}).First(&address)

	response := responses.NewAddress(address, versionURL)

	responses.SendSuccess(ctx, http.StatusOK, "retrieve-address", response, nil)
}

func RetrieveAddressList(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	addresses := []models.Address{}
	if err := database.DB.Find(&addresses).Error; err != nil {
		// TODO: Implements in future
		return
	}

	addressListResponse := []responses.Address{}
	for _, address := range addresses {
		addressListResponse = append(
			addressListResponse,
			responses.NewAddress(address, versionURL),
		)
	}

	addressList := responses.HATEOASAddressList{
		Addresses: &addressListResponse,
	}

	addressListLinks := responses.HATEOASAddressListLinks{
		Self:            responses.HREFObject{HREF: fmt.Sprintf("%s/addresses", versionURL)},
		CreateAddresses: responses.HREFObject{HREF: fmt.Sprintf("%s/addresses", versionURL)},
	}

	templateParams := []responses.HATEOASTemplateParams{
		{
			Name:          "create-addresses",
			ResourceURL:   fmt.Sprintf("%s/addresses", versionURL),
			HTTPMethod:    http.MethodPost,
			RequestStruct: requests.Address{},
		},
		{
			Name:          "create-addresses-cinemas",
			ResourceURL:   fmt.Sprintf("%s/addresses/:addressId/cinemas", versionURL),
			HTTPMethod:    http.MethodPost,
			RequestStruct: requests.Cinema{},
		},
		{
			Name:        "retrieve-cinema-list",
			ResourceURL: fmt.Sprintf("%s/addresses/:addressId/cinemas", versionURL),
			HTTPMethod:  http.MethodGet,
		},
	}
	templateJSON, err := templateFactory(versionURL, templateParams)
	if err != nil {
		// TODO: Implements in future
		return
	}

	result := responses.HATEOASResult{
		Embedded:  addressList,
		Links:     addressListLinks,
		Templates: templateJSON,
	}

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"retrieve-address-list",
		result,
		responses.HALHeaders,
	)
}
