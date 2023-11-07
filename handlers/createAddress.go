package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	model, _ := models.NewAddress(
		request.Country,
		request.State,
		request.Telephone,
		request.Description,
		request.PostalCode,
		request.Name,
	)

	if err := database.DB.Create(&model).Error; err != nil {
		responses.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	responses.SendSuccess(ctx, http.StatusOK, "CreateAddress", model)
}

func RetrieveAddress(ctx *gin.Context) {
	id := ctx.Query("id")

	model := models.Address{}
	database.DB.First(&model, id)

	responses.SendSuccess(ctx, http.StatusOK, "RetrieveAddress", model)
}
