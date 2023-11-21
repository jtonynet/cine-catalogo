package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"

	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/internal/handlers/requests"
	"github.com/jtonynet/cine-catalogo/internal/handlers/responses"
	"github.com/jtonynet/cine-catalogo/internal/hateoas"
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

	var addresses []models.Address
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

		addresses = append(addresses, address)
	}

	if err := database.DB.Create(&addresses).Error; err != nil {
		responses.SendError(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	result, err := getAddresListResult(addresses, versionURL)
	if err != nil {
		// TODO: Implements in future
		return
	}

	responses.SendSuccess(
		ctx, http.StatusOK,
		"create-addresses",
		result,
		responses.HALHeaders,
	)
}

func RetrieveAddress(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	addressId := ctx.Param("addressId")
	if !IsValidUUID(addressId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing addressId", nil)
		return
	}
	addressUUID := uuid.MustParse(addressId)

	address := models.Address{UUID: addressUUID}
	database.DB.Where(&models.Address{UUID: addressUUID}).First(&address)

	templateParams := []hateoas.TemplateParams{
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
	templateJSON, err := hateoas.TemplateFactory(versionURL, templateParams)
	if err != nil {
		// TODO: Implements in future
		return
	}

	response := responses.NewAddress(
		address,
		versionURL,
		responses.WithAddressTemplates(templateJSON),
	)

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"retrieve-address",
		response,
		nil,
	)
}

func RetrieveAddressList(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	addresses := []models.Address{}
	if err := database.DB.Find(&addresses).Error; err != nil {
		// TODO: Implements in future
		return
	}

	result, err := getAddresListResult(addresses, versionURL)
	if err != nil {
		// TODO: Implements in future
		return
	}

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"retrieve-address-list",
		result,
		responses.HALHeaders,
	)
}

func getAddresListResult(addresses []models.Address, versionURL string) (*responses.HATEOASListResult, error) {
	addressListResponse := []responses.Address{}

	for _, address := range addresses {
		addressListResponse = append(
			addressListResponse,
			*responses.NewAddress(address, versionURL),
		)
	}

	addressList := responses.HATEOASAddressList{
		Addresses: &addressListResponse,
	}

	addressListLinks := responses.HATEOASAddressListLinks{
		Self:            responses.HATEOASLink{HREF: fmt.Sprintf("%s/addresses", versionURL)},
		CreateAddresses: responses.HATEOASLink{HREF: fmt.Sprintf("%s/addresses", versionURL)},
	}

	templateParams := []hateoas.TemplateParams{
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
	templateJSON, err := hateoas.TemplateFactory(versionURL, templateParams)
	if err != nil {
		// TODO: Implements in future
		return nil, err
	}

	result := responses.HATEOASListResult{
		Embedded:  addressList,
		Links:     addressListLinks,
		Templates: templateJSON,
	}

	return &result, nil
}
