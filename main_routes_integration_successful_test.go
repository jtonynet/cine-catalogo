package main_routes_integration_test

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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
	"github.com/jtonynet/cine-catalogo/internal/handlers/requests"
	"github.com/jtonynet/cine-catalogo/internal/handlers/responses"
	"github.com/jtonynet/cine-catalogo/internal/middlewares"
	"github.com/jtonynet/cine-catalogo/internal/models"
)

type IntegrationSuccesful struct {
	suite.Suite

	cfg *config.Config

	versionURL string

	router   *gin.Engine
	routesV1 *gin.RouterGroup

	addressUUID uuid.UUID
	cinemaUUID  uuid.UUID
	movieUUID   uuid.UUID
	posterUUID  uuid.UUID

	addressResponse responses.Address

	uploadMoviePosterPath string
}

func (suite *IntegrationSuccesful) SetupSuite() {

	suite.cfg = setupConfig()
	suite.versionURL = fmt.Sprintf("%s/%s", suite.cfg.API.Host, "v1")
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)

	handlers.Init()
	database.Init(suite.cfg.Database)

	suite.addressUUID, _ = uuid.Parse("9aa904a0-feed-4502-ace8-bf9dd0e23fb5") // uuid.New()
	suite.cinemaUUID, _ = uuid.Parse("51276e29-940d-4d21-aa74-c0c4d3c5d632")  // uuid.New()
	suite.movieUUID, _ = uuid.Parse("44adac31-5290-44bf-b330-ebffe60ae0be")   // uuid.New()
	suite.posterUUID, _ = uuid.Parse("16462dd9-a701-430d-a443-4667b3a4614f")  // uuid.New()

	suite.uploadMoviePosterPath = fmt.Sprintf("%s/%s", suite.cfg.API.PostersDir, suite.movieUUID.String())
}

func (suite *IntegrationSuccesful) TearDownSuite() {
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

	database.DB.Exec(query)

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

func (suite *IntegrationSuccesful) TestV1HappyPathIntegrationSuccessful() {
	// ADDRESSES CONTEXT
	suite.addressesRoutes()

	// CINEMAS CONTEXT
	suite.cinemasRoutes()

	// MOVIES CONTEXT
	suite.moviesRoutes()

	// POSTERS CONTEXT
	suite.postersRoutes()

	// DELETES
	suite.deleteCinemaRoute()
	suite.deleteAddressRoute()
}

func (suite *IntegrationSuccesful) addressesRoutes() {
	// Create Addresses
	suite.routesV1.POST("/addresses", handlers.CreateAddresses)

	addressCreate := requests.Address{
		UUID:        suite.addressUUID,
		Country:     "BR",
		State:       "SP",
		Telephone:   "(11)0000-0000",
		Description: "Jardins Shoppings um dos mais belos de SP",
		PostalCode:  "1139050",
		Name:        "Jardins Shoppings",
	}

	addressCreateJson, err := json.Marshal(addressCreate)
	assert.NoError(suite.T(), err)
	reqAddressCreate, err := http.NewRequest("POST", "/v1/addresses", bytes.NewBuffer(addressCreateJson))
	assert.NoError(suite.T(), err)
	respAddressCreate := httptest.NewRecorder()
	suite.router.ServeHTTP(respAddressCreate, reqAddressCreate)
	assert.Equal(suite.T(), http.StatusCreated, respAddressCreate.Code)

	// Retrieve Address
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/addresses/:address_id", handlers.RetrieveAddress)

	addressUUIDRoute := fmt.Sprintf("/v1/addresses/%s", suite.addressUUID.String())

	reqAddressRetrieve, err := http.NewRequest("GET", addressUUIDRoute, nil)
	assert.NoError(suite.T(), err)
	respCinemaRetrieve := httptest.NewRecorder()
	suite.router.ServeHTTP(respCinemaRetrieve, reqAddressRetrieve)

	bodyAddressRetrieveJson := respCinemaRetrieve.Body.String()
	assert.Equal(suite.T(), http.StatusOK, respCinemaRetrieve.Code)
	assert.Equal(suite.T(), respCinemaRetrieve.Header().Get("Content-Type"), responses.JSONDefaultHeaders["Content-Type"])

	assert.Equal(suite.T(), gjson.Get(bodyAddressRetrieveJson, "uuid").String(), suite.addressUUID.String())
	assert.Equal(suite.T(), gjson.Get(bodyAddressRetrieveJson, "country").String(), addressCreate.Country)
	assert.Equal(suite.T(), gjson.Get(bodyAddressRetrieveJson, "state").String(), addressCreate.State)
	assert.Equal(suite.T(), gjson.Get(bodyAddressRetrieveJson, "telephone").String(), addressCreate.Telephone)
	assert.Equal(suite.T(), gjson.Get(bodyAddressRetrieveJson, "description").String(), addressCreate.Description)
	assert.Equal(suite.T(), gjson.Get(bodyAddressRetrieveJson, "postalCode").String(), addressCreate.PostalCode)
	assert.Equal(suite.T(), gjson.Get(bodyAddressRetrieveJson, "name").String(), addressCreate.Name)

	// Update Address
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.PATCH("/addresses/:address_id", handlers.UpdateAddress)

	addressUpdateRequest := requests.UpdateAddress{
		Telephone: "1111-1111",
	}

	addressUpdateJson, err := json.Marshal(addressUpdateRequest)
	assert.NoError(suite.T(), err)
	reqAddressUpdate, err := http.NewRequest("PATCH", addressUUIDRoute, bytes.NewBuffer(addressUpdateJson))
	assert.NoError(suite.T(), err)
	respAddressUpdate := httptest.NewRecorder()
	suite.router.ServeHTTP(respAddressUpdate, reqAddressUpdate)

	bodyAddressUpdateJson := respAddressUpdate.Body.String()
	assert.Equal(suite.T(), http.StatusOK, respAddressUpdate.Code)
	assert.Equal(suite.T(), respAddressUpdate.Header().Get("Content-Type"), responses.JSONDefaultHeaders["Content-Type"])
	assert.Equal(suite.T(), gjson.Get(bodyAddressUpdateJson, "telephone").String(), addressUpdateRequest.Telephone)

	// Retrieve Address List
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/addresses", handlers.RetrieveAddressList)

	reqRetrieveAddressList, err := http.NewRequest("GET", "/v1/addresses", nil)
	assert.NoError(suite.T(), err)
	respRetrieveAddressList := httptest.NewRecorder()
	suite.router.ServeHTTP(respRetrieveAddressList, reqRetrieveAddressList)

	bodyRetrieveAddressListJson := respRetrieveAddressList.Body.String()

	addressModel, err := models.NewAddress(
		addressCreate.UUID,
		addressCreate.Country,
		addressCreate.State,
		addressUpdateRequest.Telephone,
		addressCreate.Description,
		addressCreate.PostalCode,
		addressCreate.Name,
	)
	assert.NoError(suite.T(), err)

	addressResponse := responses.NewAddress(
		addressModel,
		suite.versionURL,
	)
	addressResponseJson, err := json.Marshal(addressResponse)
	assert.NoError(suite.T(), err)

	suite.addressResponse = addressResponse

	assert.Equal(suite.T(), http.StatusOK, respRetrieveAddressList.Code)
	assert.Contains(suite.T(), gjson.Get(bodyRetrieveAddressListJson, "_embedded.addresses").String(), string(addressResponseJson))
}

func (suite *IntegrationSuccesful) cinemasRoutes() {
	// Create Cinemas
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.POST("/addresses/:address_id/cinemas", handlers.CreateCinemas)

	cinemaCreate := requests.Cinema{
		UUID:        suite.cinemaUUID,
		Name:        "Sala Majestic IMAX 1",
		Description: "Sala IMAX com profundidade de audio",
		Capacity:    120,
	}

	cinemaCreateJson, err := json.Marshal(cinemaCreate)
	assert.NoError(suite.T(), err)
	addressUUIDCinemaRoute := fmt.Sprintf("/v1/addresses/%s/cinemas", suite.addressUUID.String())
	reqCinemasCreate, err := http.NewRequest("POST", addressUUIDCinemaRoute, bytes.NewBuffer(cinemaCreateJson))
	assert.NoError(suite.T(), err)
	respCinemasCreate := httptest.NewRecorder()

	suite.router.ServeHTTP(respCinemasCreate, reqCinemasCreate)

	assert.Equal(suite.T(), http.StatusCreated, respCinemasCreate.Code)
	assert.Equal(suite.T(), respCinemasCreate.Header().Get("Content-Type"), responses.HALHeaders["Content-Type"])

	// Retrieve Cinema
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/cinemas/:cinema_id", handlers.RetrieveCinema)

	cinemaUUIDRoute := fmt.Sprintf("/v1/cinemas/%v", suite.cinemaUUID.String())

	reqCinemaRetrieve, err := http.NewRequest("GET", cinemaUUIDRoute, nil)
	assert.NoError(suite.T(), err)
	respCinemaRetrieve := httptest.NewRecorder()
	suite.router.ServeHTTP(respCinemaRetrieve, reqCinemaRetrieve)

	bodyRetrieveCinemaJson := respCinemaRetrieve.Body.String()
	assert.Equal(suite.T(), http.StatusOK, respCinemaRetrieve.Code)
	assert.Equal(suite.T(), respCinemaRetrieve.Header().Get("Content-Type"), responses.JSONDefaultHeaders["Content-Type"])

	assert.Equal(suite.T(), gjson.Get(bodyRetrieveCinemaJson, "uuid").String(), suite.cinemaUUID.String())
	assert.Equal(suite.T(), gjson.Get(bodyRetrieveCinemaJson, "name").String(), cinemaCreate.Name)
	assert.Equal(suite.T(), gjson.Get(bodyRetrieveCinemaJson, "description").String(), cinemaCreate.Description)
	assert.Equal(suite.T(), (gjson.Get(bodyRetrieveCinemaJson, "capacity").Int()), cinemaCreate.Capacity)

	// Update Cinema
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.PATCH("/cinemas/:cinema_id", handlers.UpdateCinema)

	cinemaUpdateRequest := requests.UpdateCinema{
		Description: "Sala IMAX com profundidade de audio Surround 5D",
		Capacity:    100,
	}

	cinemaUpdateJson, err := json.Marshal(cinemaUpdateRequest)
	assert.NoError(suite.T(), err)
	reqCinemaUpdate, err := http.NewRequest("PATCH", cinemaUUIDRoute, bytes.NewBuffer(cinemaUpdateJson))
	assert.NoError(suite.T(), err)
	respCinemaUpdate := httptest.NewRecorder()
	suite.router.ServeHTTP(respCinemaUpdate, reqCinemaUpdate)

	bodyCinemaUpdateJson := respCinemaUpdate.Body.String()
	assert.Equal(suite.T(), http.StatusOK, respCinemaUpdate.Code)
	assert.Equal(suite.T(), respCinemaUpdate.Header().Get("Content-Type"), responses.JSONDefaultHeaders["Content-Type"])
	assert.Equal(suite.T(), gjson.Get(bodyCinemaUpdateJson, "description").String(), cinemaUpdateRequest.Description)
	assert.Equal(suite.T(), gjson.Get(bodyCinemaUpdateJson, "capacity").Int(), cinemaUpdateRequest.Capacity)

	// Retrieve Cinema List
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/addresses/:address_id/cinemas", handlers.RetrieveCinemaList)

	addressCinemasListUUIDRoute := fmt.Sprintf("/v1/addresses/%s/cinemas", suite.addressUUID.String())
	reqRetrieveCinemaList, err := http.NewRequest("GET", addressCinemasListUUIDRoute, nil)
	assert.NoError(suite.T(), err)
	respRetrieveCinemaList := httptest.NewRecorder()
	suite.router.ServeHTTP(respRetrieveCinemaList, reqRetrieveCinemaList)

	bodyRetrieveCinemaListJson := respRetrieveCinemaList.Body.String()

	addressCinemaListModel := models.Address{}
	err = database.DB.Where(&models.Address{UUID: suite.addressUUID}).First(&addressCinemaListModel).Error
	assert.NoError(suite.T(), err)

	cinemaModel, err := models.NewCinema(
		suite.cinemaUUID,
		addressCinemaListModel.ID,
		cinemaCreate.Name,
		cinemaUpdateRequest.Description,
		cinemaUpdateRequest.Capacity,
	)
	assert.NoError(suite.T(), err)

	cinemaResponse := responses.NewCinema(
		cinemaModel,
		suite.addressResponse.Links.Self.HREF,
		suite.versionURL,
	)
	cinemaResponseJson, err := json.Marshal(cinemaResponse)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusOK, respRetrieveCinemaList.Code)
	assert.Contains(suite.T(), gjson.Get(bodyRetrieveCinemaListJson, "_embedded.cinemas").String(), string(cinemaResponseJson))
}

func (suite *IntegrationSuccesful) moviesRoutes() {
	// Create Movies
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.POST("/movies", handlers.CreateMovies)

	ageRating := int64(14)
	published := true
	subtitled := false
	movieCreate := requests.Movie{
		UUID:        suite.movieUUID,
		Name:        "Back To The recursion 2",
		Description: "Uma aventura no tempo usando técnicas avançadas de desenvolvimento de software",
		AgeRating:   &ageRating,
		Published:   &published,
		Subtitled:   &subtitled,
	}
	movieCreateJson, err := json.Marshal(movieCreate)
	assert.NoError(suite.T(), err)
	reqMoviesCreate, err := http.NewRequest("POST", "/v1/movies", bytes.NewBuffer(movieCreateJson))
	assert.NoError(suite.T(), err)
	respMoviesCreate := httptest.NewRecorder()
	suite.router.ServeHTTP(respMoviesCreate, reqMoviesCreate)

	assert.Equal(suite.T(), http.StatusCreated, respMoviesCreate.Code)
	assert.Equal(suite.T(), respMoviesCreate.Header().Get("Content-Type"), responses.HALHeaders["Content-Type"])

	// Retrieve Movie
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/movies/:movie_id", handlers.RetrieveMovie)

	movieUUIDRoute := fmt.Sprintf("/v1/movies/%v", suite.movieUUID.String())

	reqMovieRetrieve, err := http.NewRequest("GET", movieUUIDRoute, nil)
	assert.NoError(suite.T(), err)
	respMovieRetrieve := httptest.NewRecorder()
	suite.router.ServeHTTP(respMovieRetrieve, reqMovieRetrieve)

	bodyRetrieveMovieJson := respMovieRetrieve.Body.String()
	assert.Equal(suite.T(), http.StatusOK, respMovieRetrieve.Code)
	assert.Equal(suite.T(), respMovieRetrieve.Header().Get("Content-Type"), responses.HALHeaders["Content-Type"])
	assert.Equal(suite.T(), gjson.Get(bodyRetrieveMovieJson, "uuid").String(), suite.movieUUID.String())
	assert.Equal(suite.T(), gjson.Get(bodyRetrieveMovieJson, "name").String(), movieCreate.Name)
	assert.Equal(suite.T(), gjson.Get(bodyRetrieveMovieJson, "description").String(), movieCreate.Description)
	assert.Equal(suite.T(), gjson.Get(bodyRetrieveMovieJson, "ageRating").Int(), *movieCreate.AgeRating)
	assert.Equal(suite.T(), gjson.Get(bodyRetrieveMovieJson, "published").Bool(), *movieCreate.Published)
	assert.Equal(suite.T(), gjson.Get(bodyRetrieveMovieJson, "subtitled").Bool(), *movieCreate.Subtitled)

	// Update Movie
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.PATCH("/movies/:movie_id", handlers.UpdateMovie)

	movieUpdateRequest := requests.UpdateMovie{
		Name: "Back To The recursion",
	}

	movieUpdateJson, err := json.Marshal(movieUpdateRequest)
	assert.NoError(suite.T(), err)

	reqMovieUpdate, err := http.NewRequest("PATCH", movieUUIDRoute, bytes.NewBuffer(movieUpdateJson))
	assert.NoError(suite.T(), err)
	respMovieUpdate := httptest.NewRecorder()
	suite.router.ServeHTTP(respMovieUpdate, reqMovieUpdate)

	bodyMovieUpdateJson := respMovieUpdate.Body.String()
	assert.Equal(suite.T(), http.StatusOK, respMovieUpdate.Code)
	assert.Equal(suite.T(), gjson.Get(bodyMovieUpdateJson, "name").String(), movieUpdateRequest.Name)

	// Retrieve Movie List
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/movies", handlers.RetrieveMovieList)

	reqRetrieveMovieList, err := http.NewRequest("GET", "/v1/movies", nil)
	assert.NoError(suite.T(), err)
	respRetrieveMovieList := httptest.NewRecorder()
	suite.router.ServeHTTP(respRetrieveMovieList, reqRetrieveMovieList)

	bodyRetrieveMovieListJson := respRetrieveMovieList.Body.String()

	movieModel, err := models.NewMovie(
		suite.movieUUID,
		movieUpdateRequest.Name,
		movieCreate.Description,
		*movieCreate.AgeRating,
		*movieCreate.Published,
		*movieCreate.Subtitled,
	)
	assert.NoError(suite.T(), err)

	movieResponse := responses.NewMovieListItem(
		movieModel,
		suite.cfg.API.Host,
		suite.versionURL,
	)
	movieResponseJson, err := json.Marshal(movieResponse)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusOK, respMovieRetrieve.Code)
	assert.Contains(suite.T(), gjson.Get(bodyRetrieveMovieListJson, "_embedded.movies").String(), string(movieResponseJson))

}

func (suite *IntegrationSuccesful) postersRoutes() {
	// Upload Movie Poster
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.POST("/movies/:movie_id/posters", handlers.UploadMoviePoster)

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
	suite.routesV1.GET("/movies/:movie_id/posters/:poster_id", handlers.RetrieveMoviePoster)

	movieUUIDPosterUUIDRoute := fmt.Sprintf("/v1/movies/%s/posters/%s", suite.movieUUID, suite.posterUUID)

	reqMoviePosterRetrieve, err := http.NewRequest("GET", movieUUIDPosterUUIDRoute, nil)
	assert.NoError(suite.T(), err)
	respMoviePosterRetrieve := httptest.NewRecorder()
	suite.router.ServeHTTP(respMoviePosterRetrieve, reqMoviePosterRetrieve)

	bodyRetrieveMoviePoster := respMoviePosterRetrieve.Body.String()
	assert.Equal(suite.T(), http.StatusOK, respMoviePosterRetrieve.Code)
	assert.Equal(suite.T(), respMoviePosterRetrieve.Header().Get("Content-Type"), responses.HALHeaders["Content-Type"])
	assert.Equal(suite.T(), gjson.Get(bodyRetrieveMoviePoster, "uuid").String(), posterMultPartFields["uuid"])
	assert.Equal(suite.T(), gjson.Get(bodyRetrieveMoviePoster, "name").String(), posterMultPartFields["name"])
	assert.Equal(suite.T(), gjson.Get(bodyRetrieveMoviePoster, "alternativeText").String(), posterMultPartFields["alternativeText"])

	// Update Poster
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.PATCH("/movies/:movie_id/posters/:poster_id", handlers.UpdateMoviePoster)

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

func (suite *IntegrationSuccesful) deleteCinemaRoute() {
	// Delete Cinema
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.DELETE("/cinemas/:cinema_id", handlers.DeleteCinema)

	cinemaUUIDRoute := fmt.Sprintf("/v1/cinemas/%s", suite.cinemaUUID.String())
	reqCinemaDelete, err := http.NewRequest("DELETE", cinemaUUIDRoute, nil)
	assert.NoError(suite.T(), err)
	respCinemaDelete := httptest.NewRecorder()
	suite.router.ServeHTTP(respCinemaDelete, reqCinemaDelete)

	assert.Equal(suite.T(), http.StatusNoContent, respCinemaDelete.Code)
}

func (suite *IntegrationSuccesful) deleteAddressRoute() {
	// Delete Address
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.DELETE("/addresses/:address_id", handlers.DeleteAddress)

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
	suite.Run(t, new(IntegrationSuccesful))
}
