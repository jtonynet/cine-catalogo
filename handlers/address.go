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

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		//log.Error("handlers:CreateAddress binding error %v", err.Error())
		responses.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

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
	id := ctx.Params.ByName("id")

	address := models.Address{}
	database.DB.First(&address, id)

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
