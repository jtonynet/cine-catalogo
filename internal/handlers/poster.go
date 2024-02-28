package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/internal/handlers/requests"
	"github.com/jtonynet/cine-catalogo/internal/handlers/responses"
	"github.com/jtonynet/cine-catalogo/internal/models"
)

type PosterHandler struct {
	*database.Database
}

func NewPosterHandler(db *database.Database) *PosterHandler {
	return &PosterHandler{
		Database: db,
	}
}

// @BasePath /v1

// @Summary Upload Movie Poster
// @Description Upload Movie Poster
// @Tags Movies Posters
// @Router /movies/{movie_id}/posters [post]
// @Param movie_id path string true "Movie UUID"
// @Accept mpfd
// @Param request formData requests.Poster true "Request formData"
// @Param file formData file true "binary poster data"
// @Produce json
// @Success 200 {object} responses.Poster
func (ph *PosterHandler) UploadMoviePoster(ctx *gin.Context) {

	// INFO: Param swaggo dont accepts complex type as `File *multipart.FileHeader`
	// into requests.Poster struct. Therefore, I defined the type in the Swaggo
	// annotations and use swaggerignore:"true" into requests.Poster `file` struct
	// field. see details in: https://github.com/swaggo/swag/issues/802

	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	movieId := ctx.Param("movie_id")
	if !IsValidUUID(movieId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing movie_id", nil)
		return
	}
	movieUUID := uuid.MustParse(movieId)

	var movie models.Movie
	if err := ph.Database.DB.Preload("Posters").Where(&models.Movie{UUID: movieUUID}).First(&movie).Error; err != nil {
		// TODO: Implements in future
		fmt.Println("dont found movie")
		return
	}

	if len(movie.Posters) > 0 {
		// TODO: Implements in future
		fmt.Printf("movie %s has been poster %s", movie.UUID.String(), movie.Posters[0].UUID.String())
		return
	}

	var requestForm requests.Poster
	if err := ctx.ShouldBindWith(&requestForm, binding.FormMultipart); err != nil {
		// TODO: Implements in future
		fmt.Println("dont bind requestForm %v", err)
		return
	}

	if !IsValidUUID(requestForm.UUID) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing movieId", nil)
		return
	}
	posterUUID := uuid.MustParse(requestForm.UUID)

	posterUploadedPath, err := ph.uploadPoster(ctx, movieUUID, posterUUID, requestForm.File)
	if err != nil {
		// TODO: Implements in future
		fmt.Printf("error on posterUploadedPath %v", err)
		return
	}

	poster, err := models.NewPoster(
		posterUUID,
		movie.ID,
		requestForm.Name,
		requestForm.File.Header.Get("Content-Type"),
		requestForm.AlternativeText,
		posterUploadedPath,
	)
	if err != nil {
		// TODO: Implements in future
		return
	}

	if err := ph.Database.DB.Create(&poster).Error; err != nil {
		// TODO: Implements in future
		return
	}

	resultMovie := responses.NewMovie(
		movie,
		cfg.Host,
		versionURL,
	)

	result := responses.NewPoster(
		poster,
		movie.UUID,
		resultMovie.Links.Self.HREF,
		cfg.Host,
		versionURL,
	)

	responses.SendSuccess(
		ctx,
		http.StatusCreated,
		"upload-movies-poster",
		result,
		responses.HALHeaders,
	)
}

// @Summary Update Movie Poster
// @Description Update Movie Poster
// @Tags Movies Posters
// @Router /movies/{movie_id}/posters/{poster_id} [patch]
// @Param movie_id path string true "Movie UUID"
// @Param poster_id path string true "Poster UUID"
// @Accept mpfd
// @Param request formData requests.UpdatePoster true "Request body for update"
// @Param file formData file true "binary poster data"
// @Produce json
// @Success 200 {object} responses.Poster
func (ph *PosterHandler) UpdateMoviePoster(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	movieId := ctx.Param("movie_id")
	posterId := ctx.Param("poster_id")

	if !IsValidUUID(movieId) || !IsValidUUID(posterId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing UUIDs", nil)
		return
	}

	movieUUID := uuid.MustParse(movieId)
	posterUUID := uuid.MustParse(posterId)

	existingPoster, err := ph.getPosterByMovieAndPosterUUID(posterUUID, movieUUID)
	if err != nil {
		// TODO: Implements in future
		fmt.Printf("existingPoster %v", err)
		return
	}

	var updateRequest requests.UpdatePoster
	if err := ctx.ShouldBindWith(&updateRequest, binding.FormMultipart); err != nil {
		// TODO: Implements in future
		fmt.Printf("updateRequest ShouldBindWith %v", err)
		responses.SendError(ctx, http.StatusBadRequest, "malformed request formData", nil)
		return
	}

	var posterUploadedPath string
	if updateRequest.File != nil {
		posterUploadedPath, err = ph.uploadPoster(ctx, movieUUID, posterUUID, updateRequest.File)
		if err != nil {
			// TODO: Implements in future
			fmt.Printf("posterUploadedPath %v", err)
			return
		}

		existingPoster.ContentType = updateRequest.File.Header.Get("Content-Type")
		existingPoster.Path = posterUploadedPath
	}

	if updateRequest.Name != "" {
		existingPoster.Name = updateRequest.Name
	}

	if updateRequest.AlternativeText != "" {
		existingPoster.AlternativeText = updateRequest.AlternativeText
	}

	if err := ph.Database.DB.Save(&existingPoster).Error; err != nil {
		// TODO: Implements in future
		fmt.Printf("database.DB.Save %v", err)
		return
	}

	resultMovie := responses.NewMovie(
		existingPoster.Movie,
		cfg.Host,
		versionURL,
	)

	result := responses.NewPoster(
		existingPoster,
		movieUUID,
		resultMovie.Links.Self.HREF,
		cfg.Host,
		versionURL,
	)

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"update-movies-poster",
		result,
		responses.HALHeaders,
	)
}

// @Summary Retrieve Movie Poster
// @Description Retrieve Movie Poster
// @Tags Movies Posters
// @Router /movies/{movie_id}/posters/{poster_id} [get]
// @Param movie_id path string true "Movie UUID"
// @Param poster_id path string true "Poster UUID"
// @Accept json
// @Produce json
// @Success 200 {object} responses.Poster
func (ph *PosterHandler) RetrieveMoviePoster(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	movieId := ctx.Param("movie_id")
	posterId := ctx.Param("poster_id")

	if !IsValidUUID(movieId) || !IsValidUUID(posterId) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing UUIDs", nil)
		return
	}

	movieUUID := uuid.MustParse(movieId)
	posterUUID := uuid.MustParse(posterId)

	existingPoster, err := ph.getPosterByMovieAndPosterUUID(posterUUID, movieUUID)
	if err != nil {
		// TODO: Implements in future
		fmt.Printf("existingPoster %v", err)
		return
	}

	resultMovie := responses.NewMovie(
		existingPoster.Movie,
		cfg.Host,
		versionURL,
	)

	result := responses.NewPoster(
		existingPoster,
		movieUUID,
		resultMovie.Links.Self.HREF,
		cfg.Host,
		versionURL,
	)

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"retrieve-movie-poster",
		result,
		responses.HALHeaders,
	)

}

func (ph *PosterHandler) uploadPoster(ctx *gin.Context, movieUUID, posterUUID uuid.UUID, file *multipart.FileHeader) (string, error) {
	cfg := ctx.MustGet("cfg").(config.API)

	// TODO: posters dirs, move to storages local ceph | S3 in future and manage by envVars
	uploadPath := fmt.Sprintf("%s/%s", cfg.PostersDir, movieUUID.String())
	err := os.Mkdir(uploadPath, os.ModeDir|0755) // 0644 7777 755
	if err != nil && !os.IsExist(err) {
		// TODO: Implements in future
		fmt.Printf("error on create poster directory %v ", err)
		return "", err
	}

	posterPath := filepath.Join(uploadPath, posterUUID.String()+filepath.Ext(file.Filename))
	if err := ctx.SaveUploadedFile(file, posterPath); err != nil {
		// TODO: Implements in future
		fmt.Printf("error on SaveUploadedFile %v", err)
		return "", err
	}

	return posterPath, nil
}

func (ph *PosterHandler) getPosterByMovieAndPosterUUID(posterUUID, movieUUID uuid.UUID) (models.Poster, error) {

	// TODO: Uggly query, move to model
	// or find better way:
	// https://gorm.io/docs/preload.html

	var existingPoster models.Poster
	if err := ph.Database.DB.
		Preload("Movie").
		Joins("INNER JOIN movies ON posters.movie_id = movies.id").
		Where("posters.uuid = ? AND movies.uuid = ?", posterUUID, movieUUID).
		Where("posters.deleted_at IS NULL").
		Order("posters.id").
		Limit(1).
		First(&existingPoster).Error; err != nil {
		// TODO: Implements in future
		fmt.Printf("existingPoster %v", err)
		return existingPoster, err
	}

	return existingPoster, nil
}
