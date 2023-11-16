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

func CreateCinemas(ctx *gin.Context) {
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
			responses.Cinema{
				UUID:        cinema.UUID,
				Name:        cinema.Name,
				Description: cinema.Description,
				Capacity:    cinema.Capacity,
			},
		)
	}

	responses.SendSuccess(ctx, http.StatusOK, "CreateCinemas", responseList, responses.HALHeaders)
}

// TODO:
// Retrieve Cinema List
// http://localhost:8080/api/cinemas?addressId={uuid} RetrieveCinemaList
// if addressId is dont informed, return 403 - forbbiden
//func RetrieveCinemaList (ctx *gin.Context){}

// 2eaee488-77f1-42df-b8c6-8828204ff9e3
func RetrieveCinemaList(ctx *gin.Context) {

	fmt.Println("TESTE")

	addressId, ok := ctx.GetQuery("addressId")
	if !ok {
		// if addressId is dont informed, return 403 - forbbiden
		fmt.Printf("403")
		return
	}

	addressUUID := uuid.MustParse(addressId)
	var address models.Address
	if err := database.DB.Find(&models.Address{UUID: addressUUID}).First(&address).Error; err != nil {
		fmt.Println("Cannot obtains address %v", err)
		return
	}

	var cinemas []models.Cinema
	if err := database.DB.Where(&models.Cinema{AddressID: address.ID}).Find(&cinemas).Error; err != nil {
		fmt.Println("Cannot obtains cinemas %v", err)
		return
	}

	var cinemasResponse []responses.Cinema
	for _, cinema := range cinemas {
		cinemasResponse = append(cinemasResponse,
			responses.Cinema{
				UUID:        cinema.UUID,
				Name:        cinema.Name,
				Description: cinema.Description,
				Capacity:    cinema.Capacity,
			})
	}

	responses.SendSuccess(ctx, http.StatusOK, "retrieve-cinema-list", cinemasResponse, nil)
}
