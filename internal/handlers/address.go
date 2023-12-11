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
	"github.com/jtonynet/cine-catalogo/internal/models"
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
	handler := "create-addresses"

	var requestList []requests.Address
	if err := ctx.ShouldBindBodyWith(&requestList, binding.JSON); err != nil {

		var singleRequest requests.Address
		if err := ctx.ShouldBindBodyWith(&singleRequest, binding.JSON); err != nil {
			log.WithError(err).
				WithField("handler", handler).
				Error("error on binding requests.Address")

			responses.SendError(ctx, http.StatusBadRequest, "Malformed request body.", nil)
			return
		}

		requestList = append(requestList, singleRequest)
	}

	var addresses []models.Address
	for _, request := range requestList {
		//TODO ADD UUID COLLISION MANAGEMENT
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
			log.WithError(err).
				WithField("handler", handler).
				Error("error on models.NewAddress")

			responses.SendError(ctx, http.StatusBadRequest, "Malformed request body.", nil)
			return
		}

		addresses = append(addresses, address)
	}

	if err := database.DB.Create(&addresses).Error; err != nil {
		log.WithError(err).
			WithField("handler", handler).
			Error("error on DB create addresses")

		responses.SendError(ctx, http.StatusInternalServerError, "Internal Server Error, please try again later.", nil)
		return
	}

	result, err := getAddresListResult(addresses, versionURL)
	if err != nil {
		log.WithError(err).
			WithField("handler", handler).
			Error("error on getAddresListResult")

		responses.SendError(ctx, http.StatusInternalServerError, "Internal Server Error, please try again later.", nil)
		return
	}

	responses.SendSuccess(
		ctx, http.StatusOK,
		handler,
		result,
		responses.HALHeaders,
	)
}

// @Summary Update Address
// @Description Update Address
// @Tags Addresses
// @Accept json
// @Produce json
// @Router /addresses/{address_id} [patch]
// @Param address_id path string true "Address UUID"
// @Param request body requests.UpdateAddress true "Request body"
// @Success 200 {object} responses.Address
func UpdateAddress(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")
	handler := "retrieve-address"

	addressId := ctx.Param("address_id")
	if !IsValidUUID(addressId) {
		log.WithField("handler", handler).
			Error("error invalid address_id")

		responses.SendError(ctx, http.StatusForbidden, "Malformed or missing address_id", nil)
		return
	}
	addressUUID := uuid.MustParse(addressId)

	address := models.Address{UUID: addressUUID}
	if err := database.DB.Where(&models.Address{UUID: addressUUID}).First(&address).Error; err != nil {
		log.WithError(err).
			WithField("handler", handler).
			Error("error on DB fetch address")

		responses.SendError(ctx, http.StatusForbidden, "Failed to fetch address", nil)
		return
	}

	var updateRequest requests.UpdateAddress
	if err := ctx.ShouldBind(&updateRequest); err != nil {
		log.WithError(err).
			WithField("handler", handler).
			Error("error on binding requests.UpdateAddress")

		responses.SendError(ctx, http.StatusBadRequest, "Malformed request body", nil)
		return
	}

	if updateRequest.Name != "" {
		address.Name = updateRequest.Name
	}

	if updateRequest.Country != "" {
		address.Country = updateRequest.Country
	}

	if updateRequest.State != "" {
		address.State = updateRequest.State
	}

	if updateRequest.Telephone != "" {
		address.Telephone = updateRequest.Telephone
	}

	if updateRequest.Description != "" {
		address.Description = updateRequest.Description
	}

	if updateRequest.PostalCode != "" {
		address.PostalCode = updateRequest.PostalCode
	}

	if err := database.DB.Save(&address).Error; err != nil {
		log.WithError(err).
			WithField("handler", handler).
			Error("error on DB update address")

		responses.SendError(ctx, http.StatusInternalServerError, "Internal Server Error, please try again later.", nil)
		return
	}

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
		log.WithError(err).
			WithField("handler", handler).
			Error("error on hateoas template to address")

		responses.SendError(ctx, http.StatusInternalServerError, "Internal Server Error, please try again later.", nil)
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
		handler,
		response,
		nil,
	)
}

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
	handler := "retrieve-address"

	addressId := ctx.Param("address_id")
	if !IsValidUUID(addressId) {
		log.WithField("handler", handler).
			Error("error invalid address_id")

		responses.SendError(ctx, http.StatusForbidden, "Malformed or missing address_id", nil)
		return
	}
	addressUUID := uuid.MustParse(addressId)

	address := models.Address{UUID: addressUUID}
	if err := database.DB.Where(&models.Address{UUID: addressUUID}).First(&address).Error; err != nil {
		log.WithError(err).
			WithField("handler", handler).
			Error("error on DB fetch address")

		responses.SendError(ctx, http.StatusNotFound, "Address Not Found", nil)
		return
	}

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
		log.WithError(err).
			WithField("handler", handler).
			Error("error on hateoas template to address")

		responses.SendError(ctx, http.StatusInternalServerError, "Internal Server Error, please try again later.", nil)
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
		handler,
		response,
		nil,
	)
}

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
	handler := "retrieve-address-list"

	addresses := []models.Address{}
	if err := database.DB.Find(&addresses).Error; err != nil {
		log.WithError(err).
			WithField("handler", handler).
			Error("error on DB fetch addresses")

		responses.SendError(ctx, http.StatusForbidden, "Failed to fetch addresses", nil)
		return
	}

	result, err := getAddresListResult(addresses, versionURL)
	if err != nil {
		log.WithError(err).
			WithField("handler", handler).
			Error("error on DB fetch addresses")

		responses.SendError(ctx, http.StatusForbidden, "Failed to fetch addresses", nil)
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

// @Summary Delete Address
// @Description Delete Address
// @Tags Addresses
// @Accept json
// @Produce json
// @Router /addresses/{address_id} [delete]
// @Param address_id path string true "Address UUID"
// @Success 204
func DeleteAddress(ctx *gin.Context) {
	handler := "delete-address"

	addressId := ctx.Param("address_id")
	if !IsValidUUID(addressId) {
		log.WithField("handler", handler).
			Error("error invalid address_id")

		responses.SendError(ctx, http.StatusForbidden, "Malformed or missing address_id", nil)
		return
	}
	addressUUID := uuid.MustParse(addressId)

	address := models.Address{UUID: addressUUID}
	if err := database.DB.Where(&models.Address{UUID: addressUUID}).First(&address).Error; err != nil {
		log.WithError(err).
			WithField("handler", handler).
			Error("error on DB fetch address")

		responses.SendError(ctx, http.StatusNotFound, "Address Not Found", nil)
		return
	}

	if result := database.DB.Delete(&address); result.Error != nil {
		log.WithError(result.Error).
			WithField("handler", handler).
			Error("error on DB delete address")

		responses.SendError(ctx, http.StatusInternalServerError, "Failed on delete address", nil)
		return
	}

	responses.SendSuccess(
		ctx,
		http.StatusNoContent,
		"delete-address",
		nil,
		nil,
	)
}

func getAddresListResult(addresses []models.Address, versionURL string) (*responses.HATEOASListResult, error) {
	addressListResponse := []responses.Address{}

	for _, address := range addresses {
		addressListResponse = append(
			addressListResponse,
			responses.NewAddress(address, versionURL),
		)
	}

	addressList := responses.HATEOASAddressList{
		Addresses: addressListResponse,
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
			ContentType:   "application/json",
			RequestStruct: requests.Address{},
		},
		{
			Name:          "update-address",
			ResourceURL:   fmt.Sprintf("%s/addresses/:addressId", versionURL),
			HTTPMethod:    http.MethodPatch,
			ContentType:   "application/json",
			RequestStruct: requests.UpdateAddress{},
		},
		{
			Name:        "delete-address",
			ResourceURL: fmt.Sprintf("%s/addresses/:addressId", versionURL),
			HTTPMethod:  http.MethodDelete,
			ContentType: "application/json",
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
		return nil, err
	}

	result := responses.HATEOASListResult{
		Embedded:  addressList,
		Links:     addressListLinks,
		Templates: templateJSON,
	}

	return &result, nil
}
