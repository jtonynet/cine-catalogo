package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"

	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/handlers/requests"
	"github.com/jtonynet/cine-catalogo/handlers/responses"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/internal/hateoas"
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

	responses.SendSuccess(ctx, http.StatusOK, "create-movies", responseList, nil)
}

func RetrieveMovie(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	movieId := ctx.Param("movieId")
	if !IsValidUUID(movieId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing movieId", nil)
		return
	}
	movieUUID := uuid.MustParse(movieId)

	movie := models.Movie{UUID: movieUUID}
	database.DB.Where(&models.Movie{UUID: movieUUID}).First(&movie)

	templateParams := []hateoas.TemplateParams{
		{
			Name:        "upload-movie-poster",
			ResourceURL: fmt.Sprintf("%s/movies", versionURL),
			HTTPMethod:  http.MethodPut,
		},
	}
	templateJSON, err := hateoas.TemplateFactory(versionURL, templateParams)
	if err != nil {
		// TODO: Implements in future
		return
	}

	response := *responses.NewMovie(
		movie,
		versionURL,
		responses.WithMovieTemplates(templateJSON),
	)

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"retrieve-address",
		response,
		nil,
	)
}

func RetrieveMovieList(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	movies := []models.Movie{}

	if err := database.DB.Find(&movies).Error; err != nil {
		// TODO: Implements in future
		return
	}

	movieListResponse := []responses.Movie{}
	for _, movie := range movies {
		movieListResponse = append(
			movieListResponse,
			*responses.NewMovie(
				movie,
				versionURL,
			),
		)
	}

	movieList := responses.HATEOASMovieList{
		Movies: &movieListResponse,
	}

	movieListLinks := responses.HATEOASMovieListLinks{
		Self:         responses.HATEOASLink{HREF: fmt.Sprintf("%s/movies", versionURL)},
		CreateMovies: responses.HATEOASLink{HREF: fmt.Sprintf("%s/movies", versionURL)},
	}

	templateParams := []hateoas.TemplateParams{
		{
			Name:        "retrieve-movie-list",
			ResourceURL: fmt.Sprintf("%s/movies", versionURL),
			HTTPMethod:  http.MethodGet,
		},
		{
			Name:          "create-movies",
			ResourceURL:   fmt.Sprintf("%s/movies", versionURL),
			HTTPMethod:    http.MethodPost,
			RequestStruct: requests.Movie{},
		},
		{
			Name:        "upload-movie-poster",
			ResourceURL: fmt.Sprintf("%s/movies", versionURL),
			HTTPMethod:  http.MethodPut,
			//RequestStruct: requests.Movie{},
		},
	}
	templateJSON, err := hateoas.TemplateFactory(versionURL, templateParams)
	if err != nil {
		// TODO: Implements in future
		return
	}

	result := responses.HATEOASListResult{
		Embedded:  movieList,
		Links:     movieListLinks,
		Templates: templateJSON,
	}

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"retrieve-movie-list",
		result,
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
