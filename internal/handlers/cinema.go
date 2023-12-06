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

	result, err := getCinemaListResult(cinemas, versionURL, addressId)
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
	database.DB.Where(&models.Cinema{UUID: cinemaUUID}).First(&cinema)

	templateParams := []hateoas.TemplateParams{
		{
			Name:        "retrieve-cinema",
			ResourceURL: fmt.Sprintf("%s/cinemas/%s", versionURL, cinemaId),
			ContentType: "application/json",
			HTTPMethod:  http.MethodGet,
		},
	}
	templateJSON, err := hateoas.TemplateFactory(versionURL, templateParams)
	if err != nil {
		// TODO: Implements in future
		return
	}

	response := responses.NewCinema(
		cinema,
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

	result, err := getCinemaListResult(cinemas, versionURL, addressId)
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

func getCinemaListResult(cinemas []models.Cinema, versionURL, addressId string) (*responses.HATEOASListResult, error) {
	var cinemaListResponse []responses.Cinema

	for _, cinema := range cinemas {
		cinemaListResponse = append(cinemaListResponse,
			*responses.NewCinema(cinema, versionURL),
		)
	}

	cinemaList := responses.HATEOASCinemaList{
		Cinemas: &cinemaListResponse,
	}

	cinemaListLinks := responses.HATEOASCinemaListLinks{
		Self: responses.HATEOASLink{HREF: fmt.Sprintf("%s/addresses/%s/cinemas", versionURL, addressId)},
	}

	templateParams := []hateoas.TemplateParams{
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
		return nil, err
	}

	result := &responses.HATEOASListResult{
		Embedded:  cinemaList,
		Links:     cinemaListLinks,
		Templates: templateJSON,
	}

	return result, nil
}
