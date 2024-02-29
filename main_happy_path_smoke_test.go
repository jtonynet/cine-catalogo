package main_happy_path_smoke_test

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/internal/handlers"
	"github.com/jtonynet/cine-catalogo/internal/handlers/responses"
	"github.com/jtonynet/cine-catalogo/internal/middlewares"
	"github.com/jtonynet/cine-catalogo/internal/models"
)

type HappyPath struct {
	suite.Suite

	cfg *config.Config

	versionURL string

	router   *gin.Engine
	routesV1 *gin.RouterGroup

	addressUUID uuid.UUID
	cinemaUUID  uuid.UUID
	movieUUID   uuid.UUID
	posterUUID  uuid.UUID

	addressHandler *handlers.AdrressHandler
	cinemaHandler  *handlers.CinemaHandler
	movieHandler   *handlers.MovieHandler
	posterHandler  *handlers.PosterHandler

	uploadMoviePosterPath string

	Database *database.Database
}

func (suite *HappyPath) SetupSuite() {

	suite.cfg = setupConfig()
	suite.versionURL = fmt.Sprintf("%s/%s", suite.cfg.API.Host, "v1")
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)

	db, err := database.NewDatabase(&suite.cfg.Database)
	if err != nil {
		panic("cannot connect to database")
	}

	suite.Database = db

	suite.addressHandler = handlers.NewAddressHandler(db)
	suite.cinemaHandler = handlers.NewCinemaHandler(db)
	suite.movieHandler = handlers.NewMovieHandler(db)
	suite.posterHandler = handlers.NewPosterHandler(db)

	suite.addressUUID, _ = uuid.Parse("9aa904a0-feed-4502-ace8-bf9dd0e23fb5") // uuid.New()
	suite.cinemaUUID, _ = uuid.Parse("51276e29-940d-4d21-aa74-c0c4d3c5d632")  // uuid.New()
	suite.movieUUID, _ = uuid.Parse("44adac31-5290-44bf-b330-ebffe60ae0be")   // uuid.New()
	suite.posterUUID, _ = uuid.Parse("16462dd9-a701-430d-a443-4667b3a4614f")  // uuid.New()

	suite.uploadMoviePosterPath = fmt.Sprintf("%s/%s", suite.cfg.API.PostersDir, suite.movieUUID.String())

}

func (suite *HappyPath) TearDownSuite() {
	query := fmt.Sprintf(`
		DELETE FROM cinemas WHERE uuid in ('%s');
		DELETE FROM addresses WHERE uuid in ('%s');
		DELETE FROM posters WHERE uuid in ('%s');
		DELETE FROM movies WHERE uuid in ('%s');`,
		suite.cinemaUUID.String(),
		suite.addressUUID.String(),
		suite.posterUUID.String(),
		suite.movieUUID.String(),
	)

	suite.Database.DB.Exec(query)

	err := os.RemoveAll(suite.uploadMoviePosterPath)
	if err != nil {
		fmt.Printf("Error on exclude movie poster: %v\n", err)
	}
}

func setupConfig() *config.Config {
	cfg := config.Config{}

	cfg.API.Host = "localhost:8080"
	cfg.API.StaticsDir = "web"
	cfg.API.PostersDir = "web/posters"
	cfg.API.MetricEnabled = false

	cfg.Database.Host = "localhost"
	cfg.Database.Port = "5432"
	cfg.Database.User = "api_user"
	cfg.Database.Pass = "api_pass"
	cfg.Database.DB = "cine_catalog_db"
	cfg.Database.MetricEnabled = false

	return &cfg
}

func setupRouterAndGroup(cfg config.API) (*gin.Engine, *gin.RouterGroup) {
	basePath := "/v1"

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.Use(middlewares.ConfigInject(cfg))

	return router, router.Group(basePath)
}

func (suite *HappyPath) TestV1IntegrationSuccessful() {
	suite.createAndRetrieveAddresses()
	suite.updateAndRetrieveAddressList()

	suite.createAndRetrieveCinemas()
	suite.updateAndRetrieveCinemaList()

	suite.createAndRetrieveMovies()
	suite.updateAndRetrieveMovieList()

	suite.createAndRetrieveAndUpdatePoster()

	suite.deleteCinema()
	suite.deleteAddress()
}

func (suite *HappyPath) createAndRetrieveAddresses() {
	// Create Addresses
	suite.routesV1.POST("/addresses", suite.addressHandler.CreateAddresses)

	addressCreate := fmt.Sprintf(
		`{
			"uuid":        "%s",
			"country":     "BR",
			"state":       "SP",
			"telephone":   "(11)0000-0000",
			"description": "Jardins Shoppings um dos mais belos de SP",
			"postalCode":  "1139050",
			"name":        "Jardins Shoppings"
		}`, suite.addressUUID)

	reqAddressCreate, err := http.NewRequest("POST", "/v1/addresses", bytes.NewBuffer([]byte(addressCreate)))
	assert.NoError(suite.T(), err)
	respAddressCreate := httptest.NewRecorder()
	suite.router.ServeHTTP(respAddressCreate, reqAddressCreate)
	assert.Equal(suite.T(), http.StatusCreated, respAddressCreate.Code)

	// Retrieve Address
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/addresses/:address_id", suite.addressHandler.RetrieveAddress)

	addressUUIDRoute := fmt.Sprintf("/v1/addresses/%s", suite.addressUUID.String())

	reqAddressRetrieve, err := http.NewRequest("GET", addressUUIDRoute, nil)
	assert.NoError(suite.T(), err)
	respCinemaRetrieve := httptest.NewRecorder()
	suite.router.ServeHTTP(respCinemaRetrieve, reqAddressRetrieve)

	bodyAddressRetrieveJson := respCinemaRetrieve.Body.String()
	assert.Equal(suite.T(), http.StatusOK, respCinemaRetrieve.Code)
	assert.Equal(suite.T(), respCinemaRetrieve.Header().Get("Content-Type"), responses.JSONDefaultHeaders["Content-Type"])

	assert.Equal(suite.T(), gjson.Get(bodyAddressRetrieveJson, "uuid").String(), suite.addressUUID.String())
}

func (suite *HappyPath) updateAndRetrieveAddressList() {
	// Update Address
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.PATCH("/addresses/:address_id", suite.addressHandler.UpdateAddress)

	addressUpdateRequest := `{
		"telephone": "1111-1111"
	}`

	addressUUIDRoute := fmt.Sprintf("/v1/addresses/%s", suite.addressUUID.String())

	reqAddressUpdate, err := http.NewRequest("PATCH", addressUUIDRoute, bytes.NewBuffer([]byte(addressUpdateRequest)))
	assert.NoError(suite.T(), err)
	respAddressUpdate := httptest.NewRecorder()
	suite.router.ServeHTTP(respAddressUpdate, reqAddressUpdate)

	assert.Equal(suite.T(), http.StatusOK, respAddressUpdate.Code)
	assert.Equal(suite.T(), respAddressUpdate.Header().Get("Content-Type"), responses.JSONDefaultHeaders["Content-Type"])

	// Retrieve Address List
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/addresses", suite.addressHandler.RetrieveAddressList)

	reqRetrieveAddressList, err := http.NewRequest("GET", "/v1/addresses", nil)
	assert.NoError(suite.T(), err)
	respRetrieveAddressList := httptest.NewRecorder()
	suite.router.ServeHTTP(respRetrieveAddressList, reqRetrieveAddressList)
	assert.Equal(suite.T(), http.StatusOK, respRetrieveAddressList.Code)
}

func (suite *HappyPath) createAndRetrieveCinemas() {
	// Create Cinemas
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.POST("/addresses/:address_id/cinemas", suite.cinemaHandler.CreateCinemas)

	cinemaCreate := fmt.Sprintf(
		`{
			"uuid":        "%s",
			"name":        "Sala Majestic IMAX 1",
			"description": "Sala IMAX com profundidade de audio",
			"capacity":    120
		}`, suite.cinemaUUID)

	addressUUIDCinemaRoute := fmt.Sprintf("/v1/addresses/%s/cinemas", suite.addressUUID.String())
	reqCinemasCreate, err := http.NewRequest("POST", addressUUIDCinemaRoute, bytes.NewBuffer([]byte(cinemaCreate)))
	assert.NoError(suite.T(), err)
	respCinemasCreate := httptest.NewRecorder()

	suite.router.ServeHTTP(respCinemasCreate, reqCinemasCreate)

	assert.Equal(suite.T(), http.StatusCreated, respCinemasCreate.Code)
	assert.Equal(suite.T(), respCinemasCreate.Header().Get("Content-Type"), responses.HALHeaders["Content-Type"])

	// Retrieve Cinema
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/cinemas/:cinema_id", suite.cinemaHandler.RetrieveCinema)

	cinemaUUIDRoute := fmt.Sprintf("/v1/cinemas/%v", suite.cinemaUUID.String())

	reqCinemaRetrieve, err := http.NewRequest("GET", cinemaUUIDRoute, nil)
	assert.NoError(suite.T(), err)
	respCinemaRetrieve := httptest.NewRecorder()
	suite.router.ServeHTTP(respCinemaRetrieve, reqCinemaRetrieve)

	assert.Equal(suite.T(), http.StatusOK, respCinemaRetrieve.Code)
	assert.Equal(suite.T(), respCinemaRetrieve.Header().Get("Content-Type"), responses.JSONDefaultHeaders["Content-Type"])
}

func (suite *HappyPath) updateAndRetrieveCinemaList() {
	// Update Cinema
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.PATCH("/cinemas/:cinema_id", suite.cinemaHandler.UpdateCinema)

	cinemaUUIDRoute := fmt.Sprintf("/v1/cinemas/%v", suite.cinemaUUID.String())

	description := "Sala IMAX com profundidade de audio Surround 5D"
	capacity := 100
	cinemaUpdateRequest := fmt.Sprintf(
		`{
			"description": "%s",
			"capacity":    %v
		}`, description, capacity)

	reqCinemaUpdate, err := http.NewRequest("PATCH", cinemaUUIDRoute, bytes.NewBuffer([]byte(cinemaUpdateRequest)))
	assert.NoError(suite.T(), err)
	respCinemaUpdate := httptest.NewRecorder()
	suite.router.ServeHTTP(respCinemaUpdate, reqCinemaUpdate)

	assert.Equal(suite.T(), http.StatusOK, respCinemaUpdate.Code)
	assert.Equal(suite.T(), respCinemaUpdate.Header().Get("Content-Type"), responses.JSONDefaultHeaders["Content-Type"])

	// Retrieve Cinema List
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/addresses/:address_id/cinemas", suite.cinemaHandler.RetrieveCinemaList)

	addressCinemasListUUIDRoute := fmt.Sprintf("/v1/addresses/%s/cinemas", suite.addressUUID.String())
	reqRetrieveCinemaList, err := http.NewRequest("GET", addressCinemasListUUIDRoute, nil)
	assert.NoError(suite.T(), err)
	respRetrieveCinemaList := httptest.NewRecorder()
	suite.router.ServeHTTP(respRetrieveCinemaList, reqRetrieveCinemaList)

	addressCinemaListModel := models.Address{}
	err = suite.Database.DB.Where(&models.Address{UUID: suite.addressUUID}).First(&addressCinemaListModel).Error
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusOK, respRetrieveCinemaList.Code)
}

func (suite *HappyPath) createAndRetrieveMovies() {
	// Create Movies
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.POST("/movies", suite.movieHandler.CreateMovies)

	movieCreate := fmt.Sprintf(
		`{
			"uuid":        "%s",
			"name":        "Back To The recursion 2",
			"description": "Uma aventura no tempo usando técnicas avançadas de desenvolvimento de software",
			"ageRating":   14,
			"published":   true,
			"subtitled":   false
		}`, suite.movieUUID)

	reqMoviesCreate, err := http.NewRequest("POST", "/v1/movies", bytes.NewBuffer([]byte(movieCreate)))
	assert.NoError(suite.T(), err)
	respMoviesCreate := httptest.NewRecorder()
	suite.router.ServeHTTP(respMoviesCreate, reqMoviesCreate)

	assert.Equal(suite.T(), http.StatusCreated, respMoviesCreate.Code)
	assert.Equal(suite.T(), respMoviesCreate.Header().Get("Content-Type"), responses.HALHeaders["Content-Type"])

	// Retrieve Movie
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/movies/:movie_id", suite.movieHandler.RetrieveMovie)

	movieUUIDRoute := fmt.Sprintf("/v1/movies/%v", suite.movieUUID.String())

	reqMovieRetrieve, err := http.NewRequest("GET", movieUUIDRoute, nil)
	assert.NoError(suite.T(), err)
	respMovieRetrieve := httptest.NewRecorder()
	suite.router.ServeHTTP(respMovieRetrieve, reqMovieRetrieve)

	assert.Equal(suite.T(), http.StatusOK, respMovieRetrieve.Code)
	assert.Equal(suite.T(), respMovieRetrieve.Header().Get("Content-Type"), responses.HALHeaders["Content-Type"])
}

func (suite *HappyPath) updateAndRetrieveMovieList() {
	// Update Movie
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.PATCH("/movies/:movie_id", suite.movieHandler.UpdateMovie)

	movieUUIDRoute := fmt.Sprintf("/v1/movies/%v", suite.movieUUID.String())

	movieUpdateRequest := `{
		"name": "Back To The recursion"
	}`

	reqMovieUpdate, err := http.NewRequest("PATCH", movieUUIDRoute, bytes.NewBuffer([]byte(movieUpdateRequest)))
	assert.NoError(suite.T(), err)
	respMovieUpdate := httptest.NewRecorder()
	suite.router.ServeHTTP(respMovieUpdate, reqMovieUpdate)

	assert.Equal(suite.T(), http.StatusOK, respMovieUpdate.Code)

	// Retrieve Movie List
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/movies", suite.movieHandler.RetrieveMovieList)

	reqRetrieveMovieList, err := http.NewRequest("GET", "/v1/movies", nil)
	assert.NoError(suite.T(), err)
	respRetrieveMovieList := httptest.NewRecorder()
	suite.router.ServeHTTP(respRetrieveMovieList, reqRetrieveMovieList)

	assert.Equal(suite.T(), http.StatusOK, respRetrieveMovieList.Code)
}

func (suite *HappyPath) createAndRetrieveAndUpdatePoster() {
	// Upload Movie Poster
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.POST("/movies/:movie_id/posters", suite.posterHandler.UploadMoviePoster)

	imageContentType := "image/png"

	posterPath := "./docs/assets/images/posters/back_to_the_recursion.png"
	posterFile, err := os.Open(posterPath)
	assert.NoError(suite.T(), err)
	defer posterFile.Close()

	fileInfo, err := posterFile.Stat()
	assert.NoError(suite.T(), err)
	fileBuffer := make([]byte, fileInfo.Size())
	posterFile.Read(fileBuffer)
	fileBytes := bytes.NewReader(fileBuffer)
	posterFileUploadedMD5 := calculateMD5(fileBuffer)

	posterMultPartRequestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(posterMultPartRequestBody)
	posterFileHeader := make(textproto.MIMEHeader)
	posterFileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filepath.Base(posterPath)))
	posterFileHeader.Set("Content-Type", imageContentType)
	posterFilePart, err := writer.CreatePart(posterFileHeader)
	assert.NoError(suite.T(), err)

	io.Copy(posterFilePart, fileBytes)

	posterMultPartFields := map[string]string{
		"uuid":            suite.posterUUID.String(),
		"name":            "Back To The Recursion 2",
		"alternativeText": "Uma aventura no tempo usando técnicas avançadas de desenvolvimento de software",
	}

	for key, value := range posterMultPartFields {
		err := writer.WriteField(key, value)
		assert.NoError(suite.T(), err)
	}
	writer.Close()

	movieUUIDPostersRoute := fmt.Sprintf("/v1/movies/%s/posters", suite.movieUUID.String())
	reqUploadPoster, err := http.NewRequest("POST", movieUUIDPostersRoute, posterMultPartRequestBody)
	assert.NoError(suite.T(), err)
	reqUploadPoster.Header.Set("Content-Type", writer.FormDataContentType())
	respUploadPoster := httptest.NewRecorder()
	suite.router.ServeHTTP(respUploadPoster, reqUploadPoster)

	assert.Equal(suite.T(), http.StatusCreated, respUploadPoster.Code)
	assert.Equal(suite.T(), respUploadPoster.Header().Get("Content-Type"), responses.HALHeaders["Content-Type"])
	assert.DirExists(suite.T(), suite.uploadMoviePosterPath)

	// Retrieve Static Poster File
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.router.Static("/web", suite.cfg.API.StaticsDir)

	movieUUIDPosterUUIDFileRoute := fmt.Sprintf("/web/posters/%s/%s.png", suite.movieUUID.String(), suite.posterUUID)
	reqRetrieveStaticPosterFile, err := http.NewRequest("GET", movieUUIDPosterUUIDFileRoute, nil)
	assert.NoError(suite.T(), err)
	respRetrieveStaticPosterFile := httptest.NewRecorder()
	suite.router.ServeHTTP(respRetrieveStaticPosterFile, reqRetrieveStaticPosterFile)

	assert.Equal(suite.T(), http.StatusOK, respRetrieveStaticPosterFile.Code)
	assert.Equal(suite.T(), imageContentType, respRetrieveStaticPosterFile.Header().Get("Content-Type"))

	posterFileRetrievedMD5 := calculateMD5(respRetrieveStaticPosterFile.Body.Bytes())
	assert.Equal(suite.T(), posterFileUploadedMD5, posterFileRetrievedMD5)

	// Retrieve Poster
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/movies/:movie_id/posters/:poster_id", suite.posterHandler.RetrieveMoviePoster)

	movieUUIDPosterUUIDRoute := fmt.Sprintf("/v1/movies/%s/posters/%s", suite.movieUUID, suite.posterUUID)

	reqMoviePosterRetrieve, err := http.NewRequest("GET", movieUUIDPosterUUIDRoute, nil)
	assert.NoError(suite.T(), err)
	respMoviePosterRetrieve := httptest.NewRecorder()
	suite.router.ServeHTTP(respMoviePosterRetrieve, reqMoviePosterRetrieve)

	assert.Equal(suite.T(), http.StatusOK, respMoviePosterRetrieve.Code)

	// Update Poster
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.PATCH("/movies/:movie_id/posters/:poster_id", suite.posterHandler.UpdateMoviePoster)

	posterUpdatePath := "./docs/assets/images/posters/back_to_the_recursion.png"
	posterUpdateFile, err := os.Open(posterUpdatePath)
	assert.NoError(suite.T(), err)
	defer posterUpdateFile.Close()

	fileUpdateInfo, err := posterUpdateFile.Stat()
	assert.NoError(suite.T(), err)
	fileUpdateBuffer := make([]byte, fileUpdateInfo.Size())
	posterUpdateFile.Read(fileUpdateBuffer)
	fileUpdateBytes := bytes.NewReader(fileUpdateBuffer)

	posterUpdateMultPartRequestBody := &bytes.Buffer{}
	writerUpdate := multipart.NewWriter(posterUpdateMultPartRequestBody)
	posterUpdateFileHeader := make(textproto.MIMEHeader)
	posterUpdateFileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filepath.Base(posterUpdatePath)))
	posterUpdateFileHeader.Set("Content-Type", imageContentType)
	posterUpdateFilePart, err := writerUpdate.CreatePart(posterUpdateFileHeader)
	assert.NoError(suite.T(), err)

	io.Copy(posterUpdateFilePart, fileUpdateBytes)

	posterUpdateMultPartFields := map[string]string{
		"name": "Back To The Recursion",
	}

	for key, value := range posterUpdateMultPartFields {
		err := writerUpdate.WriteField(key, value)
		assert.NoError(suite.T(), err)
	}
	writerUpdate.Close()

	movieUUIDPostersUUIDRoute := fmt.Sprintf("/v1/movies/%s/posters/%s", suite.movieUUID.String(), suite.posterUUID.String())
	reqUpdatePoster, err := http.NewRequest("PATCH", movieUUIDPostersUUIDRoute, posterUpdateMultPartRequestBody)
	assert.NoError(suite.T(), err)
	reqUpdatePoster.Header.Set("Content-Type", writerUpdate.FormDataContentType())
	respUpdatePoster := httptest.NewRecorder()
	suite.router.ServeHTTP(respUpdatePoster, reqUpdatePoster)

	assert.Equal(suite.T(), http.StatusOK, respUpdatePoster.Code)
	assert.Equal(suite.T(), respUpdatePoster.Header().Get("Content-Type"), responses.HALHeaders["Content-Type"])
}

func (suite *HappyPath) deleteCinema() {
	// Delete Cinema
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.DELETE("/cinemas/:cinema_id", suite.cinemaHandler.DeleteCinema)

	cinemaUUIDRoute := fmt.Sprintf("/v1/cinemas/%s", suite.cinemaUUID.String())
	reqCinemaDelete, err := http.NewRequest("DELETE", cinemaUUIDRoute, nil)
	assert.NoError(suite.T(), err)
	respCinemaDelete := httptest.NewRecorder()
	suite.router.ServeHTTP(respCinemaDelete, reqCinemaDelete)

	assert.Equal(suite.T(), http.StatusNoContent, respCinemaDelete.Code)
}

func (suite *HappyPath) deleteAddress() {
	// Delete Address
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.DELETE("/addresses/:address_id", suite.addressHandler.DeleteAddress)

	addressUUIDRoute := fmt.Sprintf("/v1/addresses/%s", suite.addressUUID.String())
	reqAddressDelete, err := http.NewRequest("DELETE", addressUUIDRoute, nil)
	assert.NoError(suite.T(), err)
	respAddressDelete := httptest.NewRecorder()
	suite.router.ServeHTTP(respAddressDelete, reqAddressDelete)

	assert.Equal(suite.T(), http.StatusNoContent, respAddressDelete.Code)
}

func calculateMD5(buffer []byte) string {
	hash := md5.New()
	_, _ = hash.Write(buffer)
	return hex.EncodeToString(hash.Sum(nil))
}

func TestIntegrationSuccessfulSuite(t *testing.T) {
	suite.Run(t, new(HappyPath))
}
