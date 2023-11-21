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

func CreateCinemas(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	addressId := ctx.Param("addressId")
	if !IsValidUUID(addressId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing addressId", nil)
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
			uuid.New(),
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

func RetrieveCinema(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	cinemaId := ctx.Param("cinemaId")
	if !IsValidUUID(cinemaId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing cinemaId", nil)
		return
	}
	cinemaUUID := uuid.MustParse(cinemaId)

	cinema := models.Cinema{UUID: cinemaUUID}
	database.DB.Where(&models.Cinema{UUID: cinemaUUID}).First(&cinema)

	templateParams := []hateoas.TemplateParams{
		{
			Name:        "retrieve-cinema",
			ResourceURL: fmt.Sprintf("%s/cinemas/%s", versionURL, cinemaId),
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

	addressId := ctx.Param("addressId")
	if !IsValidUUID(addressId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing addressId", nil)
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
			ResourceURL: fmt.Sprintf("%s/addresses/:addressId/cinemas", versionURL),
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
