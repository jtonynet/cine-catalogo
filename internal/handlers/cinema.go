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

type CinemaHandler struct {
	*database.Database
}

func NewCinemaHandler(db *database.Database) *CinemaHandler {
	return &CinemaHandler{
		Database: db,
	}
}

// @BasePath /v1

// @Summary Create Addresses Cinemas
// @Description Create List of Cinemas
// @Tags Addresses Cinemas
// @Accept json
// @Produce json
// @Router /addresses/{address_id}/cinemas [post]
// @Param address_id path string true "Address UUID"
// @Param request body []requests.Cinema true "Request body"
// @Success 200 {object} responses.HATEOASListResult
func (ch *CinemaHandler) CreateCinemas(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")
	handler := "create-cinemas"

	addressId := ctx.Param("address_id")
	if !IsValidUUID(addressId) {

		log.WithField("origin", handler).
			Error("error invalid address_id")

		responses.SendError(ctx, http.StatusForbidden, "Malformed or missing address_id", nil)
		return
	}
	addressUUID := uuid.MustParse(addressId)

	var address models.Address
	if err := ch.Database.DB.Where(&models.Address{UUID: addressUUID}).First(&address).Error; err != nil {

		log.WithError(err).
			WithField("origin", handler).
			Error("error on DB fetch address")

		responses.SendError(ctx, http.StatusNotFound, "Address Not Found", nil)
		return
	}

	var requestList []requests.Cinema
	if err := ctx.ShouldBindBodyWith(&requestList, binding.JSON); err != nil {

		var singleRequest requests.Cinema
		if err := ctx.ShouldBindBodyWith(&singleRequest, binding.JSON); err != nil {

			log.WithError(err).
				WithField("origin", handler).
				Error("error on binding requests.Cinema")

			responses.SendError(ctx, http.StatusBadRequest, "Malformed request body.", nil)
			return
		}

		requestList = append(requestList, singleRequest)
	}

	var cinemas []models.Cinema
	for _, request := range requestList {
		//TODO ADD UUID COLLISION MANAGEMENT
		cinema, err := models.NewCinema(
			request.UUID,
			address.ID,
			request.Name,
			request.Description,
			request.Capacity,
		)
		if err != nil {
			log.WithError(err).
				WithField("origin", handler).
				Error("error on models.NewCinema")

			responses.SendError(ctx, http.StatusBadRequest, "Malformed request body.", nil)
			return
		}

		cinemas = append(cinemas, cinema)
	}

	if err := ch.Database.DB.Create(&cinemas).Error; err != nil {
		log.WithError(err).
			WithField("origin", handler).
			Error("error on DB create cinemas")

		responses.SendError(ctx, http.StatusInternalServerError, "Internal Server Error, please try again later.", nil)
		return
	}

	result, err := ch.getCinemaListResult(cinemas, address, versionURL)
	if err != nil {
		log.WithError(err).
			WithField("origin", handler).
			Error("error on getCinemaListResult")

		responses.SendError(ctx, http.StatusInternalServerError, "Internal Server Error, please try again later.", nil)
		return
	}

	responses.SendSuccess(
		ctx,
		http.StatusCreated,
		handler,
		result,
		responses.HALHeaders,
	)
}

// @Summary Retrieve Cinema
// @Description Retrieve one Cinema
// @Tags Cinemas
// @Produce json
// @Router /cinemas/{cinema_id} [get]
// @Param cinema_id path string true "Cinema UUID"
// @Success 200 {object} responses.Cinema
func (ch *CinemaHandler) RetrieveCinema(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")
	handler := "retrieve-cinema"

	cinemaId := ctx.Param("cinema_id")
	if !IsValidUUID(cinemaId) {

		log.WithField("origin", handler).
			Error("error invalid cinema_id")

		responses.SendError(ctx, http.StatusInternalServerError, "Malformed or missing address_id", nil)
		return
	}
	cinemaUUID := uuid.MustParse(cinemaId)

	cinema := models.Cinema{UUID: cinemaUUID}
	if err := ch.Database.DB.Preload("Address").Where(&models.Cinema{UUID: cinemaUUID}).First(&cinema).Error; err != nil {
		log.WithError(err).
			WithField("origin", handler).
			Error("error on DB fetch cinema")

		responses.SendError(ctx, http.StatusNotFound, "Cinema Not Found", nil)
		return
	}

	templateJSON, err := ch.getCinemasTemplates(versionURL)
	if err != nil {
		log.WithError(err).
			WithField("origin", handler).
			Error("error on hateoas template to cinema")

		responses.SendError(ctx, http.StatusInternalServerError, "Internal Server Error, please try again later.", nil)
		return
	}

	addressResponse := responses.NewAddress(cinema.Address, versionURL)

	response := responses.NewCinema(
		cinema,
		addressResponse.Links.Self.HREF,
		versionURL,
		responses.WithCinemaTemplates(templateJSON),
	)

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		handler,
		response,
		nil,
	)
}

// @Summary Delete Cinema
// @Description Delete Cinema
// @Tags Cinemas
// @Accept json
// @Produce json
// @Router /cinemas/{cinema_id} [delete]
// @Param cinema_id path string true "Cinema UUID"
// @Success 204
func (ch *CinemaHandler) DeleteCinema(ctx *gin.Context) {
	handler := "delete-cinema"

	cinemaId := ctx.Param("cinema_id")
	if !IsValidUUID(cinemaId) {
		log.WithField("handler", handler).
			Error("error invalid cinema_id")

		responses.SendError(ctx, http.StatusForbidden, "Malformed or missing address_id", nil)
		return
	}
	cinemaUUID := uuid.MustParse(cinemaId)

	cinema := models.Cinema{UUID: cinemaUUID}
	if err := ch.Database.DB.Where(&models.Cinema{UUID: cinemaUUID}).First(&cinema).Error; err != nil {
		log.WithError(err).
			WithField("origin", handler).
			Error("error on DB fetch cinema")

		responses.SendError(ctx, http.StatusNotFound, "Cinema Not Found", nil)
		return
	}

	if result := ch.Database.DB.Delete(&cinema); result.Error != nil {
		log.WithError(result.Error).
			WithField("origin", handler).
			Error("error on DB delete cinema")

		responses.SendError(ctx, http.StatusInternalServerError, "Failed on delete address", nil)
		return
	}

	responses.SendSuccess(
		ctx,
		http.StatusNoContent,
		handler,
		nil,
		nil,
	)
}

// @Summary Update Cinema
// @Description Update Cinema
// @Tags Cinemas
// @Accept json
// @Produce json
// @Router /cinemas/{cinema_id} [patch]
// @Param cinema_id path string true "Cinema UUID"
// @Param request body requests.UpdateCinema true "Request body"
// @Success 200 {object} responses.Cinema
func (ch *CinemaHandler) UpdateCinema(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")
	handler := "retrieve-cinema"

	cinemaId := ctx.Param("cinema_id")
	if !IsValidUUID(cinemaId) {
		log.WithField("origin", handler).
			Error("error invalid cinema_id")

		responses.SendError(ctx, http.StatusForbidden, "Malformed or missing cinema_id", nil)
		return
	}
	cinemaUUID := uuid.MustParse(cinemaId)

	cinema := models.Cinema{UUID: cinemaUUID}
	if err := ch.Database.DB.Preload("Address").Where(&models.Cinema{UUID: cinemaUUID}).First(&cinema).Error; err != nil {
		log.WithError(err).
			WithField("origin", handler).
			Error("error on DB fetch cinema")

		responses.SendError(ctx, http.StatusForbidden, "Failed to fetch cinema", nil)
		return
	}

	var updateRequest requests.UpdateCinema
	if err := ctx.ShouldBindBodyWith(&updateRequest, binding.JSON); err != nil {
		log.WithError(err).
			WithField("origin", handler).
			Error("error on binding requests.UpdateCinema")

		responses.SendError(ctx, http.StatusBadRequest, "Malformed request body", nil)
		return
	}

	if updateRequest.Name != "" {
		cinema.Name = updateRequest.Name
	}

	if updateRequest.Description != "" {
		cinema.Description = updateRequest.Description
	}

	if updateRequest.Capacity > 0 {
		cinema.Capacity = updateRequest.Capacity
	}

	if err := ch.Database.DB.Save(&cinema).Error; err != nil {
		log.WithError(err).
			WithField("origin", handler).
			Error("error on DB update cinema")

		responses.SendError(ctx, http.StatusInternalServerError, "Internal Server Error, please try again later.", nil)
		return
	}

	templateJSON, err := ch.getCinemasTemplates(versionURL)
	if err != nil {
		log.WithError(err).
			WithField("origin", handler).
			Error("error on hateoas template to cinema")

		responses.SendError(ctx, http.StatusInternalServerError, "Internal Server Error, please try again later.", nil)
		return
	}

	addressResponse := responses.NewAddress(cinema.Address, versionURL)

	response := responses.NewCinema(
		cinema,
		addressResponse.Links.Self.HREF,
		versionURL,
		responses.WithCinemaTemplates(templateJSON),
	)

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		handler,
		response,
		nil,
	)
}

// @Summary Retrieve Cinema List
// @Description Retrieve List all Cinemas from one Address
// @Tags Addresses Cinemas
// @Produce json
// @Success 200 {object} responses.HATEOASListResult
// @Router /addresses/{address_id}/cinemas [get]
// @Param address_id path string true "Address UUID"
func (ch *CinemaHandler) RetrieveCinemaList(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")
	handler := "retrieve-cinema-list"

	addressId := ctx.Param("address_id")
	if !IsValidUUID(addressId) {
		log.WithField("origin", handler).
			Error("error invalid address_id")

		responses.SendError(ctx, http.StatusForbidden, "Malformed or missing address_id", nil)
		return
	}
	addressUUID := uuid.MustParse(addressId)

	address := models.Address{UUID: addressUUID}
	if err := ch.Database.DB.Where(&models.Address{UUID: addressUUID}).First(&address).Error; err != nil {
		log.WithError(err).
			WithField("origin", handler).
			Error("error on DB fetch address")

		responses.SendError(ctx, http.StatusNotFound, "Address Not Found", nil)
		return
	}

	cinemas := []models.Cinema{}
	if err := ch.Database.DB.Where(&models.Cinema{AddressID: address.ID}).Find(&cinemas).Error; err != nil {
		log.WithError(err).
			WithField("origin", handler).
			Error("error on DB fetch cinemas")

		responses.SendError(ctx, http.StatusInternalServerError, "Internal Server Error, please try again later.", nil)
		return
	}

	result, err := ch.getCinemaListResult(cinemas, address, versionURL)
	if err != nil {
		log.WithError(err).
			WithField("origin", handler).
			Error("error on getCinemaListResult")

		responses.SendError(ctx, http.StatusInternalServerError, "Internal Server Error, please try again later.", nil)
		return
	}

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		handler,
		result,
		responses.HALHeaders,
	)
}

func (ch *CinemaHandler) getCinemaListResult(cinemas []models.Cinema, address models.Address, versionURL string) (*responses.HATEOASListResult, error) {
	var cinemaListResponse []responses.Cinema

	addressResponse := responses.NewAddress(address, versionURL)

	for _, cinema := range cinemas {
		cinemaListResponse = append(cinemaListResponse,
			responses.NewCinema(
				cinema,
				addressResponse.Links.Self.HREF,
				versionURL,
			),
		)
	}

	cinemaList := responses.HATEOASCinemaList{
		Cinemas: cinemaListResponse,
	}

	cinemaListLinks := responses.HATEOASCinemaListLinks{
		Self: responses.HATEOASLink{HREF: fmt.Sprintf("%s/addresses/%s/cinemas", versionURL, address.UUID)},
	}

	templateJSON, err := ch.getCinemasTemplates(versionURL)
	if err != nil {
		// TODO: Implements in future
		return nil, err
	}

	result := &responses.HATEOASListResult{
		Embedded:  cinemaList,
		Links:     cinemaListLinks,
		Templates: templateJSON,
	}

	return result, nil
}

func (ch *CinemaHandler) getCinemasTemplates(
	versionURL string,
) (interface{}, error) {
	templateParams := []hateoas.TemplateParams{
		{
			Name:          "update-cinema",
			ResourceURL:   fmt.Sprintf("%s/cinemas/:address_id", versionURL),
			ContentType:   "application/json",
			HTTPMethod:    http.MethodPatch,
			RequestStruct: requests.UpdateCinema{},
		},
		{
			Name:          "delete-cinema",
			ResourceURL:   fmt.Sprintf("%s/cinemas/:address_id", versionURL),
			ContentType:   "application/json",
			HTTPMethod:    http.MethodDelete,
			RequestStruct: requests.UpdateCinema{},
		},
	}
	templateJSON, err := hateoas.TemplateFactory(versionURL, templateParams)
	if err != nil {
		return nil, err
	}

	return templateJSON, nil
}
