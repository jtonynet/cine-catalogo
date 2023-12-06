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

// @BasePath /v1

// @Summary Create Addresses
// @Description Create List of Addresses
// @Tags Addresses
// @Accept json
// @Produce json
// @Param request body []requests.Address true "Request body"
// @Success 200 {object} responses.HATEOASListResult
// @Router /addresses [post]
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
			request.UUID,
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

// @BasePath /v1

// @Summary Retrieve Address
// @Description Retrieve one Address
// @Tags Addresses
// @Accept json
// @Produce json
// @Param address_id path string true "UUID of the address"
// @Success 200 {object} responses.Address
// @Router /addresses/{address_id} [get]
func RetrieveAddress(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	addressId := ctx.Param("address_id")
	if !IsValidUUID(addressId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing address_id", nil)
		return
	}
	addressUUID := uuid.MustParse(addressId)

	address := models.Address{UUID: addressUUID}
	database.DB.Where(&models.Address{UUID: addressUUID}).First(&address)

	templateParams := []hateoas.TemplateParams{
		{
			Name:          "create-addresses-cinemas",
			ResourceURL:   fmt.Sprintf("%s/addresses/:address_id/cinemas", versionURL),
			HTTPMethod:    http.MethodPost,
			ContentType:   "application/json",
			RequestStruct: requests.Cinema{},
		},
		{
			Name:        "retrieve-cinema-list",
			ResourceURL: fmt.Sprintf("%s/addresses/:address_id/cinemas", versionURL),
			ContentType: "application/json",
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

// @BasePath /v1

// @Summary Retrieve Address List
// @Description Retrieve List all Address
// @Tags Addresses
// @Accept json
// @Produce json
// @Success 200 {object} responses.HATEOASListResult
// @Router /addresses [get]
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
		Self:                   responses.HATEOASLink{HREF: fmt.Sprintf("%s/addresses", versionURL)},
		CreateAddresses:        responses.HATEOASLink{HREF: fmt.Sprintf("%s/addresses", versionURL)},
		CreateAddressesCinemas: responses.HATEOASLink{HREF: fmt.Sprintf("%s/addresses/%s/cinemas", versionURL)},
	}

	templateParams := []hateoas.TemplateParams{
		{
			Name:          "create-addresses",
			ResourceURL:   fmt.Sprintf("%s/addresses", versionURL),
			HTTPMethod:    http.MethodPost,
			ContentType:   "application/json",
			RequestStruct: requests.Address{},
		},
		{
			Name:          "create-addresses-cinemas",
			ResourceURL:   fmt.Sprintf("%s/addresses/:addressId/cinemas", versionURL),
			HTTPMethod:    http.MethodPost,
			ContentType:   "application/json",
			RequestStruct: requests.Cinema{},
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
