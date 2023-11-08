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

func CreateCinema(ctx *gin.Context) {
	request := requests.Cinema{}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		//TODO: Implements in future
		return
	}

	cinema, err := models.NewCinema(
		uuid.New(),
		request.Name,
		request.Description,
		request.Capacity,
	)
	if err != nil {
		//TODO: Implements in future
		return
	}

	if err := database.DB.Create(&cinema).Error; err != nil {
		//TODO: Implements in future
		return
	}

	response := responses.Cinema{
		UUID:        cinema.UUID,
		Name:        cinema.Name,
		Description: cinema.Description,
		Capacity:    cinema.Capacity,
	}

	responses.SendSuccess(ctx, http.StatusOK, "CreateCinema", response)
}
