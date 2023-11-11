package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"

	"github.com/jtonynet/cine-catalogo/handlers/requests"
	"github.com/jtonynet/cine-catalogo/handlers/responses"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/models"
)

func CreateMovies(ctx *gin.Context) {
	// INFO: Accepts requestList || singleRequest (transforms into a list with a single value)
	var requestList []requests.Movie
	if err := ctx.ShouldBindBodyWith(&requestList, binding.JSON); err != nil {

		var singleRequest requests.Movie
		if err := ctx.ShouldBindBodyWith(&singleRequest, binding.JSON); err != nil {
			// TODO: Implements in future
			return
		}

		requestList = append(requestList, singleRequest)
	}

	movieList := []models.Movie{}
	for _, request := range requestList {
		movie, err := models.NewMovie(
			uuid.New(),
			request.Name,
			request.Description,
			request.AgeRating,
			*request.Subtitled,
			request.Poster,
		)
		if err != nil {
			// TODO: Implements in future
			return
		}

		movieList = append(movieList, movie)
	}

	if err := database.DB.Create(&movieList).Error; err != nil {
		// TODO: Implements in future
		return
	}

	responseList := []responses.Movie{}
	for _, movie := range movieList {
		responseList = append(responseList,
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

	responses.SendSuccess(ctx, http.StatusOK, "CreateMovies", responseList)
}

func RetrieveMovieList(ctx *gin.Context) {
	movies := []models.Movie{}

	if err := database.DB.Find(&movies).Error; err != nil {
		// TODO: Implements in future
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
