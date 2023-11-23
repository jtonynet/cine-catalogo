package handlers

import (
	"fmt"
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
	"github.com/jtonynet/cine-catalogo/models"
)

// @BasePath /v1

// @Summary Upload Movie Poster
// @Description Upload Movie Poster
// @Tags Movies
// @Router /movies/{movie_id}/posters [post]
// @Param movie_id path string true "Movie UUID"
// @Accept mpfd
// @Param request formData requests.Poster true "Request formData"
// @Param file formData file true "binary poster data"
// @Produce json
// @Success 200 {object} responses.Poster
func UploadMoviePoster(ctx *gin.Context) {

	// INFO: Param swaggo dont accepts complex type as `File *multipart.FileHeader`
	// into requests.Poster struct. Therefore, I defined the type in the Swaggo
	// annotations and use swaggerignore:"true" into requests.Poster `file` struct
	// field. see details in: https://github.com/swaggo/swag/issues/802

	cfg := ctx.MustGet("cfg").(config.API)
	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	movieUUID := uuid.MustParse(ctx.Param("movie_id"))

	var movie models.Movie
	if err := database.DB.Where(&models.Movie{UUID: movieUUID}).First(&movie).Error; err != nil {
		// TODO: Implements in future
		fmt.Println("dont found movie")
		return
	}

	var requestForm requests.Poster
	if err := ctx.ShouldBindWith(&requestForm, binding.FormMultipart); err != nil {
		// TODO: Implements in future
		fmt.Println("dont bind requestForm %v", err)
		return
	}

	// TODO: posters dirs, move to ceph s3 in future and manage by envVars
	uploadPath := fmt.Sprintf("%s/%s", cfg.PostersDir, movieUUID.String())
	err := os.Mkdir(uploadPath, 0644) //0644 7777
	if err != nil && !os.IsExist(err) {
		// TODO: Implements in future
		fmt.Println("error on create poster directory %v", err)
		return
	}

	if !IsValidUUID(requestForm.UUID) {
		responses.SendError(ctx, http.StatusForbidden, "malformed or missing movieId", nil)
		return
	}
	posterUUID := uuid.MustParse(requestForm.UUID)

	posterPath := filepath.Join(uploadPath, posterUUID.String()+filepath.Ext(requestForm.File.Filename))
	if err := ctx.SaveUploadedFile(requestForm.File, posterPath); err != nil {
		// TODO: Implements in future
		return
	}

	poster, err := models.NewPoster(
		posterUUID,
		movie.ID,
		requestForm.Name,
		requestForm.File.Header.Get("Content-Type"),
		requestForm.AlternativeText,
		posterPath,
	)
	if err != nil {
		// TODO: Implements in future
		return
	}

	if err := database.DB.Create(&poster).Error; err != nil {
		// TODO: Implements in future
		return
	}

	result := responses.NewPoster(
		poster,
		movie.UUID,
		cfg.Host,
		versionURL,
		nil,
	)

	responses.SendSuccess(
		ctx,
		http.StatusOK,
		"upload-movies-poster",
		result,
		responses.HALHeaders,
	)
}
