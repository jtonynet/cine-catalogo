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

// @Summary Create Addresses Cinemas
// @Description Create List of Cinemas
// @Tags Addresses Cinemas
// @Accept json
// @Produce json
// @Router /addresses/{address_id}/cinemas [post]
// @Param address_id path string true "Address UUID"
// @Param request body []requests.Cinema true "Request body"
// @Success 200 {object} responses.HATEOASListResult
func CreateCinemas(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	addressId := ctx.Param("address_id")
	if !IsValidUUID(addressId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing address_id", nil)
		return
	}
	addressUUID := uuid.MustParse(addressId)

	var address models.Address
	if err := database.DB.Where(&models.Address{UUID: addressUUID}).First(&address).Error; err != nil {
		// TODO: Implements in future
		return
	}

	var requestList []requests.Cinema
	if err := ctx.ShouldBindBodyWith(&requestList, binding.JSON); err != nil {
		var singleRequest requests.Cinema
		if err := ctx.ShouldBindBodyWith(&singleRequest, binding.JSON); err != nil {
			// TODO: Implements in future
			return
		}

		requestList = append(requestList, singleRequest)
	}

	var cinemas []models.Cinema
	for _, request := range requestList {
		cinema, err := models.NewCinema(
			request.UUID,
			address.ID,
			request.Name,
			request.Description,
			request.Capacity,
		)
		if err != nil {
			// TODO: Implements in future
			return
		}

		cinemas = append(cinemas, cinema)
	}

	if err := database.DB.Create(&cinemas).Error; err != nil {
		// TODO: Implements in future
		return
	}

	result, err := getCinemaListResult(cinemas, address, versionURL)
	if err != nil {
		// TODO: Implements in future
		return
	}

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"create-cinemas",
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
func RetrieveCinema(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	cinemaId := ctx.Param("cinema_id")
	if !IsValidUUID(cinemaId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing cinema_id", nil)
		return
	}
	cinemaUUID := uuid.MustParse(cinemaId)

	cinema := models.Cinema{UUID: cinemaUUID}
	if err := database.DB.Preload("Address").Where(&models.Cinema{UUID: cinemaUUID}).First(&cinema).Error; err != nil {
		responses.SendError(ctx, http.StatusNotFound, "cinema not found", nil)
		return
	}

	templateJSON, err := getCinemasTemplates(versionURL)
	if err != nil {
		// TODO: Implements in future
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
		"retrieve-cinema",
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
func DeleteCinema(ctx *gin.Context) {

	cinemaId := ctx.Param("cinema_id")
	if !IsValidUUID(cinemaId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing cinema_id", nil)
		return
	}
	cinemaUUID := uuid.MustParse(cinemaId)

	cinema := models.Cinema{UUID: cinemaUUID}
	if err := database.DB.Where(&models.Cinema{UUID: cinemaUUID}).First(&cinema).Error; err != nil {
		responses.SendError(ctx, http.StatusNotFound, "cinema not found", nil)
		return
	}

	if result := database.DB.Delete(&cinema); result.Error != nil {
		responses.SendError(ctx, http.StatusInternalServerError, "failed to delete cinema", nil)
		return
	}

	responses.SendSuccess(
		ctx,
		http.StatusNoContent,
		"delete-cinema",
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
func UpdateCinema(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	cinemaId := ctx.Param("cinema_id")
	if !IsValidUUID(cinemaId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing cinema_id", nil)
		return
	}
	cinemaUUID := uuid.MustParse(cinemaId)

	cinema := models.Cinema{UUID: cinemaUUID}
	if err := database.DB.Preload("Address").Where(&models.Cinema{UUID: cinemaUUID}).First(&cinema).Error; err != nil {
		responses.SendError(ctx, http.StatusForbidden, "dont fetch cinema", nil)
		return
	}

	var updateRequest requests.UpdateCinema
	if err := ctx.ShouldBind(&updateRequest); err != nil {
		// TODO: Implements in future
		fmt.Printf("updateRequest ShouldBindJSON %v", err)
		responses.SendError(ctx, http.StatusBadRequest, "malformed request body", nil)
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

	if err := database.DB.Save(&cinema).Error; err != nil {
		// TODO: Implements in future
		fmt.Printf("database.DB.Save %v", err)
		return
	}

	templateJSON, err := getCinemasTemplates(versionURL)
	if err != nil {
		// TODO: Implements in future
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
		"retrieve-cinema",
		response,
		nil,
	)
}

// @Summary Retrieve Cinema List
// @Description Retrieve List all Cinemas from one Address
// @Tags Addresses Cinemas
// @Produce json
// @Success 200 {object} responses.MovieListResult
// @Router /addresses/{address_id}/cinemas [get]
// @Param address_id path string true "Address UUID"
func RetrieveCinemaList(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	addressId := ctx.Param("address_id")
	if !IsValidUUID(addressId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing address_id", nil)
		return
	}
	addressUUID := uuid.MustParse(addressId)

	address := models.Address{UUID: addressUUID}
	if err := database.DB.Where(&models.Address{UUID: addressUUID}).First(&address).Error; err != nil {
		fmt.Println("Cannot obtains address %v", err)
		return
	}

	cinemas := []models.Cinema{}
	if err := database.DB.Where(&models.Cinema{AddressID: address.ID}).Find(&cinemas).Error; err != nil {
		// TODO: Implements in future
		return
	}

	result, err := getCinemaListResult(cinemas, address, versionURL)
	if err != nil {
		// TODO: Implements in future
		return
	}

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"retrieve-cinema-list",
		result,
		responses.HALHeaders,
	)
}

func getCinemaListResult(cinemas []models.Cinema, address models.Address, versionURL string) (*responses.HATEOASListResult, error) {
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

	templateJSON, err := getCinemasTemplates(versionURL)
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

func getCinemasTemplates(
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
		// TODO: Implements in future
		return nil, err
	}

	return templateJSON, nil
}
