package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"

	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/internal/handlers/requests"
	"github.com/jtonynet/cine-catalogo/internal/handlers/responses"
	"github.com/jtonynet/cine-catalogo/internal/hateoas"
	"github.com/jtonynet/cine-catalogo/internal/models"
)

// @BasePath /v1

// @Summary Create Movies
// @Description Create List of Movies
// @Tags Movies
// @Accept json
// @Produce json
// @Param request body []requests.Movie true "Request body"
// @Success 200 {object} responses.HATEOASListResult
// @Router /movies [post]
func CreateMovies(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")
	handler := "create-movies"

	var requestList []requests.Movie
	if err := ctx.ShouldBindBodyWith(&requestList, binding.JSON); err != nil {
		var singleRequest requests.Movie
		if err := ctx.ShouldBindBodyWith(&singleRequest, binding.JSON); err != nil {
			// TODO: Implements in future
			return
		}

		requestList = append(requestList, singleRequest)
	}

	movies := []models.Movie{}
	for _, request := range requestList {

		movie, err := models.NewMovie(
			request.UUID,
			request.Name,
			request.Description,
			*request.AgeRating,
			*request.Published,
			*request.Subtitled,
		)
		if err != nil {
			// TODO: Implements in future
			return
		}

		movies = append(movies, movie)
	}

	if err := database.DB.Create(&movies).Error; err != nil {
		// TODO: Implements in future
		return
	}

	result, err := getMovieListResult(movies, cfg.Host, versionURL)
	if err != nil {
		// TODO: Implements in future
		return
	}

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		handler,
		result,
		responses.HALHeaders,
	)
}

// @Summary Retrieve Movie
// @Description Retrieve one Movie
// @Tags Movies
// @Accept json
// @Produce json
// @Param movie_id path string true "UUID of the movie"
// @Success 200 {object} responses.Movie
// @Router /movies/{movie_id} [get]
func RetrieveMovie(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	movieId := ctx.Param("movie_id")
	if !IsValidUUID(movieId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing movie_id", nil)
		return
	}
	movieUUID := uuid.MustParse(movieId)

	movie := models.Movie{UUID: movieUUID}
	if err := database.DB.Preload("Posters").Where(&models.Movie{UUID: movieUUID}).First(&movie).Error; err != nil {
		responses.SendError(ctx, http.StatusForbidden, "dont fetch Movie and Poster", nil)
		return
	}

	templateJSON, err := getMoviesTemplates(versionURL)
	if err != nil {
		// TODO: Implements in future
		return
	}

	response := responses.NewMovie(
		movie,
		cfg.Host,
		versionURL,
		responses.WithMovieTemplates(templateJSON),
	)

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"retrieve-movie",
		response,
		responses.HALHeaders,
	)
}

// @Summary Update Movie
// @Description Update Movie
// @Tags Movies
// @Router /movies/{movie_id} [patch]
// @Param movie_id path string true "Movie UUID"
// @Accept json
// @Param request body requests.UpdateMovie true "Request body for update"
// @Produce json
// @Success 200 {object} responses.Movie
func UpdateMovie(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	movieId := ctx.Param("movie_id")
	if !IsValidUUID(movieId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing movie_id", nil)
		return
	}
	movieUUID := uuid.MustParse(movieId)

	movie := models.Movie{UUID: movieUUID}
	if err := database.DB.Preload("Posters").Where(&models.Movie{UUID: movieUUID}).First(&movie).Error; err != nil {
		responses.SendError(ctx, http.StatusForbidden, "dont fetch Movie and Poster", nil)
		return
	}

	var updateRequest requests.UpdateMovie
	if err := ctx.ShouldBind(&updateRequest); err != nil {
		// TODO: Implements in future
		fmt.Printf("updateRequest ShouldBindJSON %v", err)
		responses.SendError(ctx, http.StatusBadRequest, "malformed request body", nil)
		return
	}

	if updateRequest.Name != "" {
		movie.Name = updateRequest.Name
	}

	if updateRequest.Description != "" {
		movie.Description = updateRequest.Description
	}

	if updateRequest.AgeRating != nil {
		movie.AgeRating = *updateRequest.AgeRating
	}

	if updateRequest.Published != nil {
		movie.Published = *updateRequest.Published
	}

	if updateRequest.Subtitled != nil {
		movie.Subtitled = *updateRequest.Subtitled
	}

	if err := database.DB.Save(&movie).Error; err != nil {
		// TODO: Implements in future
		fmt.Printf("database.DB.Save %v", err)
		return
	}

	templateJSON, err := getMoviesTemplates(versionURL)
	if err != nil {
		// TODO: Implements in future
		return
	}

	response := responses.NewMovie(
		movie,
		cfg.Host,
		versionURL,
		responses.WithMovieTemplates(templateJSON),
	)

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"update-movie",
		response,
		responses.HALHeaders,
	)
}

// @Summary Retrieve Movie List
// @Description Retrieve List all Movies
// @Tags Movies
// @Accept json
// @Produce json
// @Success 200 {object} responses.HATEOASListResult
// @Router /movies [get]
func RetrieveMovieList(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	movies := []models.Movie{}
	if err := database.DB.Preload("Posters").Find(&movies).Error; err != nil {
		// TODO: Implements in future
		return
	}

	result, err := getMovieListResult(movies, cfg.Host, versionURL)
	if err != nil {
		// TODO: Implements in future
		return
	}

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"retrieve-movie-list",
		result,
		responses.HALHeaders,
	)
}

func getMovieListResult(movies []models.Movie, baseURL, versionURL string) (*responses.HATEOASListResult, error) {
	movieListResponse := []responses.MovieListItem{}
	posterListResponse := []responses.Poster{}

	for _, movie := range movies {
		m := responses.NewMovieListItem(
			movie,
			baseURL,
			versionURL,
		)

		movieListResponse = append(
			movieListResponse,
			m,
		)

		if len(movie.Posters) > 0 {
			posterListResponse = append(
				posterListResponse,
				responses.NewPoster(
					movie.Posters[0],
					movie.UUID,
					m.Links.Self.HREF,
					baseURL,
					versionURL,
				),
			)
		}
	}

	movieAndPosterList := responses.HATEOASMovieAndPostersList{
		Movies:  &movieListResponse,
		Posters: &posterListResponse,
	}

	movieListLinks := responses.HATEOASMovieListLinks{
		Self:         responses.HATEOASLink{HREF: fmt.Sprintf("%s/movies", versionURL)},
		CreateMovies: responses.HATEOASLink{HREF: fmt.Sprintf("%s/movies", versionURL)},
	}

	templateJSON, err := getMoviesTemplates(versionURL)
	if err != nil {
		// TODO: Implements in future
		return nil, err
	}

	result := responses.HATEOASListResult{
		Embedded:  movieAndPosterList,
		Links:     movieListLinks,
		Templates: templateJSON,
	}

	return &result, nil
}

func getMoviesTemplates(
	versionURL string,
) (interface{}, error) {
	templateParams := []hateoas.TemplateParams{
		{
			Name:          "create-movies",
			ResourceURL:   fmt.Sprintf("%s/movies", versionURL),
			HTTPMethod:    http.MethodPost,
			ContentType:   "application/json",
			RequestStruct: requests.Movie{},
		},
		{
			Name:          "update-movie",
			ResourceURL:   fmt.Sprintf("%s/movies/:movie_id", versionURL),
			HTTPMethod:    http.MethodPatch,
			ContentType:   "application/json",
			RequestStruct: requests.Movie{},
		},
		{
			Name:          "upload-movie-poster",
			ResourceURL:   fmt.Sprintf("%s/movies/:movie_id/posters", versionURL),
			HTTPMethod:    http.MethodPost,
			ContentType:   "multipart/form-data",
			RequestStruct: requests.Poster{},
		},
		{
			Name:          "update-movie-poster",
			ResourceURL:   fmt.Sprintf("%s/movies/:movie_id/posters", versionURL),
			HTTPMethod:    http.MethodPatch,
			ContentType:   "multipart/form-data",
			RequestStruct: requests.Poster{},
		},
	}
	templateJSON, err := hateoas.TemplateFactory(versionURL, templateParams)
	if err != nil {
		// TODO: Implements in future
		return nil, err
	}

	return templateJSON, nil
}
