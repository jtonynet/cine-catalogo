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

func CreateCinemas(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	addressUUID := uuid.MustParse(ctx.Param("addressId"))

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

	cinemaList := []models.Cinema{}
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

		cinemaList = append(cinemaList, cinema)
	}

	if err := database.DB.Create(&cinemaList).Error; err != nil {
		// TODO: Implements in future
		return
	}

	responseList := []responses.Cinema{}
	for _, cinema := range cinemaList {
		responseList = append(responseList,
			responses.NewCinema(cinema, versionURL),
		)
	}

	responses.SendSuccess(ctx, http.StatusOK, "CreateCinemas", responseList, responses.HALHeaders)
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

	response := responses.NewCinema(cinema, versionURL)

	responses.SendSuccess(ctx, http.StatusOK, "retrieve-cinema", response, nil)
}

func RetrieveCinemaList(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	addressId := ctx.Param("addressId")
	if !IsValidUUID(addressId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed addressId", nil)
		return
	}

	var address models.Address
	addressUUID := uuid.MustParse(addressId)
	if err := database.DB.Find(&models.Address{UUID: addressUUID}).First(&address).Error; err != nil {
		fmt.Println("Cannot obtains address %v", err)
		return
	}

	var cinemas []models.Cinema
	if err := database.DB.Where(&models.Cinema{AddressID: address.ID}).Find(&cinemas).Error; err != nil {
		fmt.Println("Cannot obtains cinemas %v", err)
		return
	}

	var cinemaListResponse []responses.Cinema
	for _, cinema := range cinemas {
		cinemaListResponse = append(cinemaListResponse,
			responses.NewCinema(cinema, versionURL),
		)
	}

	selfURL := fmt.Sprintf("%s/addresses/%s/cinemas", versionURL, addressId)
	cinemaListLinks := responses.HATEOASCinemaListLinks{
		Self: responses.HREFObject{HREF: selfURL},
	}

	cinemaList := responses.HATEOASCinemaList{
		Cinemas: &cinemaListResponse,
	}

	templateJSON, err := getHATEOASCinemaTemplate(selfURL, versionURL)
	if err != nil {
		// TODO: Implements in future
		return
	}

	result := responses.HATEOASResult{
		Embedded:  cinemaList,
		Links:     cinemaListLinks,
		Templates: templateJSON,
	}

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"retrieve-cinema-list",
		result,
		responses.HALHeaders,
	)
}

func getHATEOASCinemaTemplate(baseURL, selfURL string) (interface{}, error) {
	root := hateoas.NewRootDocument(baseURL)

	retrieveCinemaListGet, err := hateoas.NewResource(
		"retrieve-cinema-list",
		selfURL,
		http.MethodGet,
	)
	if err != nil {
		// TODO: implements on future
		return nil, err
	}
	root.AddResource(retrieveCinemaListGet)

	rootEncoded, err := root.Encode()
	if err != nil {
		// TODO: implements on future
		return nil, err
	}

	templateString := gjson.Get(string(rootEncoded), "_templates").String()
	var templateJSON interface{}
	json.Unmarshal([]byte(templateString), &templateJSON)

	return templateJSON, nil
}
