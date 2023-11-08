package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/handlers/requests"
	"github.com/jtonynet/cine-catalogo/handlers/responses"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/models"
)

func CreateAddress(ctx *gin.Context) {
	request := requests.CreateAddress{}

	ctx.ShouldBindJSON(&request)

	address, _ := models.NewAddress(
		uuid.New(),
		request.Country,
		request.State,
		request.Telephone,
		request.Description,
		request.PostalCode,
		request.Name,
	)

	if err := database.DB.Create(&address).Error; err != nil {
		responses.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response := responses.Address{
		UUID:        address.UUID,
		Country:     address.Country,
		State:       address.State,
		Telephone:   address.Telephone,
		Description: address.Description,
		PostalCode:  address.PostalCode,
		Name:        address.Name,
	}

	responses.SendSuccess(ctx, http.StatusOK, "CreateAddress", response)
}

func RetrieveAddress(ctx *gin.Context) {
	uuid := uuid.MustParse(ctx.Params.ByName("uuid"))

	address := models.Address{UUID: uuid}
	database.DB.Where(&models.Address{UUID: uuid}).First(&address)

	response := responses.Address{
		UUID:        address.UUID,
		Country:     address.Country,
		State:       address.State,
		Telephone:   address.Telephone,
		Description: address.Description,
		PostalCode:  address.PostalCode,
		Name:        address.Name,
	}

	responses.SendSuccess(ctx, http.StatusOK, "RetrieveAddress", response)
}

func RetrieveAddressList(ctx *gin.Context) {
	addresses := []models.Address{}

	if err := database.DB.Find(&addresses).Error; err != nil {
		//TODO: Implements in future
		return
	}

	response := []responses.Address{}
	for _, address := range addresses {
		response = append(
			response,
			responses.Address{
				UUID:        address.UUID,
				Country:     address.Country,
				State:       address.State,
				Telephone:   address.Telephone,
				Description: address.Description,
				PostalCode:  address.PostalCode,
				Name:        address.Name,
			},
		)
	}

	responses.SendSuccess(ctx, http.StatusOK, "RetrieveAddressList", response)
}
