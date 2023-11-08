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

func CreateMovie(ctx *gin.Context) {
	request := requests.Movie{}

	err := ctx.ShouldBind(&request)
	if err != nil {
		//TODO: Implements in future
		return
	}

	movie, err := models.NewMovie(
		uuid.New(),
		request.Name,
		request.Description,
		request.AgeRating,
		*request.Subtitled,
		request.Poster,
	)
	if err != nil {
		//TODO: Implements in future
		return
	}

	if err := database.DB.Create(&movie).Error; err != nil {
		//TODO: Implements in future
		return
	}

	response := responses.Movie{
		UUID:        movie.UUID,
		Name:        movie.Name,
		Description: movie.Description,
		AgeRating:   movie.AgeRating,
		Subtitled:   movie.Subtitled,
		Poster:      movie.Poster,
	}

	responses.SendSuccess(ctx, http.StatusOK, "CreateMovie", response)
}

func RetrieveMovieList(ctx *gin.Context) {
	movies := []models.Movie{}

	if err := database.DB.Find(&movies).Error; err != nil {
		//TODO: Implements in future
		return
	}

	response := []responses.Movie{}
	for _, movie := range movies {
		response = append(
			response,
			responses.Movie{
				UUID:        movie.UUID,
				Name:        movie.Name,
				Description: movie.Description,
				AgeRating:   movie.AgeRating,
				Subtitled:   movie.Subtitled,
				Poster:      movie.Poster,
			},
		)
	}

	responses.SendSuccess(ctx, http.StatusOK, "RetrieveMovieList", response)
}
