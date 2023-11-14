package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"

	"github.com/jtonynet/cine-catalogo/handlers/requests"
	"github.com/jtonynet/cine-catalogo/handlers/responses"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/models"
)

func CreateMovies(ctx *gin.Context) {
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

	responses.SendSuccess(ctx, http.StatusOK, "CreateMovies", responseList, nil)
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

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"RetrieveMovieList",
		response,
		responses.HALHeaders,
	)
}

func UploadMoviePoster(ctx *gin.Context) {
	movieUUID := uuid.MustParse(ctx.Param("movieId"))

	var movie models.Movie
	if err := database.DB.Where(&models.Movie{UUID: movieUUID}).First(&movie).Error; err != nil {
		// TODO: Implements in future
		fmt.Println("dont found movie")
		return
	}

	file, err := ctx.FormFile("poster")
	if err != nil {
		// TODO: Implements in future
		fmt.Printf("dont sended poster file %v", err)
		return
	}

	// TODO: posters dirs, move to ceph s3 in future and manage by envVars
	uploadPath := "./posters/"
	posterPath := filepath.Join(uploadPath, movieUUID.String()+filepath.Ext(file.Filename))
	if err := ctx.SaveUploadedFile(file, posterPath); err != nil {
		// TODO: Implements in future
		fmt.Println("dont save poster file")
		return
	}

	movie.Poster = posterPath
	if err := database.DB.Save(&movie).Error; err != nil {
		// TODO: Implements in future
		fmt.Println("dont update movie")
		return
	}
}
