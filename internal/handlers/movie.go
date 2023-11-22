package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/internal/handlers/requests"
	"github.com/jtonynet/cine-catalogo/internal/handlers/responses"
	"github.com/jtonynet/cine-catalogo/internal/hateoas"
	"github.com/jtonynet/cine-catalogo/models"
)

func CreateMovies(ctx *gin.Context) {
	// https://gin-gonic.com/docs/examples/upload-file/multiple-file/
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	form, _ := ctx.MultipartForm()

	names, _ := form.Value["name[]"]
	descriptions, _ := form.Value["description[]"]
	ageRatings, _ := form.Value["age_rating[]"]
	subtitleds, _ := form.Value["subtitled[]"]

	posters, ok := form.File["poster[]"]
	if !ok {
		// TODO: Implements in future
		fmt.Printf("dont get sended poster file")
		return
	}

	uploadPath := cfg.PostersDir
	movies := []models.Movie{}
	for idx, poster := range posters {
		movieUUID := uuid.New()
		posterPath := filepath.Join(uploadPath, movieUUID.String()+filepath.Ext(poster.Filename))
		ctx.SaveUploadedFile(poster, posterPath)

		ageRating, _ := strconv.ParseInt(ageRatings[idx], 10, 64)
		subtitled, _ := strconv.ParseBool(subtitleds[idx])

		movie, err := models.NewMovie(
			movieUUID,
			names[idx],
			descriptions[idx],
			ageRating,
			subtitled,
			posterPath,
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
		"retrieve-movie-list",
		result,
		responses.HALHeaders,
	)
}

// @BasePath /v1

// @Summary Retrieve Movie
// @Description Retrieve one Movie
// @Tags Movies
// @Accept json
// @Produce json
// @Param movieId path string true "UUID of the movie"
// @Success 200 {object} responses.Movie
// @Router /movies/{movieId} [get]
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
			Name:        "update-movie",
			ResourceURL: fmt.Sprintf("%s/movies/:movieId", versionURL),
			HTTPMethod:  http.MethodPatch,
		},
	}
	templateJSON, err := hateoas.TemplateFactory(versionURL, templateParams)
	if err != nil {
		// TODO: Implements in future
		return
	}

	response := responses.NewMovie(
		movie,
		templateJSON,
		cfg.Host,
		versionURL,
	)

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"retrieve-address",
		response,
		nil,
	)
}

// @BasePath /v1

// @Summary Retrieve List Movies
// @Description Retrieve List all Movies
// @Tags Movies
// @Accept json
// @Produce json
// @Success 200 {object} responses.MovieListResult
// @Router /movies [get]
func RetrieveMovieList(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	movies := []models.Movie{}

	if err := database.DB.Find(&movies).Error; err != nil {
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

func getMovieListResult(movies []models.Movie, baseURL, versionURL string) (*responses.MovieListResult, error) {
	movieListResponse := []responses.MovieListItem{}
	posterListResponse := []responses.HATEOASPosterLinks{}

	for _, movie := range movies {
		movieListResponse = append(
			movieListResponse,
			responses.NewMovieListItem(
				movie,
				baseURL,
				versionURL,
			),
		)

		posterListResponse = append(
			posterListResponse,
			*responses.NewPosterLinks(
				baseURL,
				movie.Poster,
			),
		)
	}

	movieAndPosterList := responses.HATEOASMovieAndPostersList{
		Movies:  &movieListResponse,
		Posters: &posterListResponse,
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
			Name:        "update-movie",
			ResourceURL: fmt.Sprintf("%s/movies/:movieId", versionURL),
			HTTPMethod:  http.MethodPatch,
			//RequestStruct: requests.Movie{},
		},
	}
	templateJSON, err := hateoas.TemplateFactory(versionURL, templateParams)
	if err != nil {
		// TODO: Implements in future
		return nil, err
	}

	result := responses.MovieListResult{
		Embedded:  movieAndPosterList,
		Links:     movieListLinks,
		Templates: templateJSON,
	}

	return &result, nil
}
