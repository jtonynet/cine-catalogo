package main_test

import (
	"bytes"
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

/*
TODO: Using gin engine or ginTestEngine for tests, context params bugfixes (move to ADR or another doc):
	https://medium.com/nerd-for-tech/testing-rest-api-in-go-with-testify-and-mockery-c31ea2cc88f9
	https://forum.golangbridge.org/t/how-to-test-gin-gonic-handler-function-within-a-function/33334/2
	https://github.com/gin-gonic/gin/issues/1292
	https://github.com/gin-gonic/gin/pull/2803
	https://github.com/gin-gonic/gin/issues/2778
	https://github.com/gin-gonic/gin/issues/2816
*/

type IntegrationSuccesfulSuite struct {
	suite.Suite
	cfg      *config.Config
	router   *gin.Engine
	routesV1 *gin.RouterGroup

	versionURL string

	addressUUID uuid.UUID
	cinemaUUID  uuid.UUID
	movieUUID   uuid.UUID
	posterUUID  uuid.UUID
}

func (suite *IntegrationSuccesfulSuite) SetupSuite() {

	suite.cfg = setupConfig()
	suite.versionURL = fmt.Sprintf("%s/%s", suite.cfg.API.Host, "v1")
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)

	database.Init(suite.cfg.Database)

	suite.addressUUID, _ = uuid.Parse("9aa904a0-feed-4502-ace8-bf9dd0e23fb5") // uuid.New()  //
	suite.cinemaUUID, _ = uuid.Parse("51276e29-940d-4d21-aa74-c0c4d3c5d632")  // uuid.New()  //
	suite.movieUUID, _ = uuid.Parse("44adac31-5290-44bf-b330-ebffe60ae0be")   // uuid.New()  //
	suite.posterUUID, _ = uuid.Parse("16462dd9-a701-430d-a443-4667b3a4614f")  // uuid.New()  //
}

func (suite *IntegrationSuccesfulSuite) TearDownSuite() {
	query := fmt.Sprintf(`
	 DELETE FROM cinemas WHERE uuid in ('%v');
	 DELETE FROM addresses WHERE uuid in ('%v');
	 DELETE FROM posters WHERE uuid in ('%v');
	 DELETE FROM movies WHERE uuid in ('%v');`,
		suite.cinemaUUID.String(),
		suite.addressUUID.String(),
		suite.posterUUID.String(),
		suite.movieUUID.String())

	database.DB.Exec(query)

	uploadPath := fmt.Sprintf("%s/%s", suite.cfg.API.PostersDir, suite.movieUUID.String())
	err := os.RemoveAll(uploadPath)
	if err != nil {
		fmt.Printf("Error on exclude movie poster: %v\n", err)
	}
}

func setupConfig() *config.Config {
	cfg := config.Config{}

	cfg.API.Host = "catalogo-api-test"
	cfg.API.StaticsDir = "web"
	cfg.API.PostersDir = "../../web/posters"
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

func (suite *IntegrationSuccesfulSuite) TestV1HappyPathIntegrationSuccessful() {

	// ADDRESSES CONTEXT

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
	addressCreateJson, _ := json.Marshal(addressCreate)
	reqAddressCreate, _ := http.NewRequest("POST", "/v1/addresses", bytes.NewBuffer(addressCreateJson))
	respAddressCreate := httptest.NewRecorder()
	suite.router.ServeHTTP(respAddressCreate, reqAddressCreate)
	assert.Equal(suite.T(), http.StatusCreated, respAddressCreate.Code)

	// Retrieve Address
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/addresses/:address_id", handlers.RetrieveAddress)

	addressUUIDRoute := fmt.Sprintf("/v1/addresses/%s", suite.addressUUID.String())

	reqAddressRetrieve, _ := http.NewRequest("GET", addressUUIDRoute, nil)
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

	addressUpdateJson, _ := json.Marshal(addressUpdateRequest)
	reqAddressUpdate, _ := http.NewRequest("PATCH", addressUUIDRoute, bytes.NewBuffer(addressUpdateJson))
	respAddressUpdate := httptest.NewRecorder()
	suite.router.ServeHTTP(respAddressUpdate, reqAddressUpdate)

	bodyAddressUpdateJson := respAddressUpdate.Body.String()
	assert.Equal(suite.T(), http.StatusOK, respAddressUpdate.Code)
	assert.Equal(suite.T(), respAddressUpdate.Header().Get("Content-Type"), responses.JSONDefaultHeaders["Content-Type"])
	assert.Equal(suite.T(), gjson.Get(bodyAddressUpdateJson, "telephone").String(), addressUpdateRequest.Telephone)

	// Retrieve Address List
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/addresses", handlers.RetrieveAddressList)

	reqRetrieveAddressList, _ := http.NewRequest("GET", "/v1/addresses", nil)
	respRetrieveAddressList := httptest.NewRecorder()
	suite.router.ServeHTTP(respRetrieveAddressList, reqRetrieveAddressList)

	bodyRetrieveAddressListJson := respRetrieveAddressList.Body.String()

	addressModel, _ := models.NewAddress(
		addressCreate.UUID,
		addressCreate.Country,
		addressCreate.State,
		addressUpdateRequest.Telephone,
		addressCreate.Description,
		addressCreate.PostalCode,
		addressCreate.Name,
	)

	addressResponse := responses.NewAddress(
		addressModel,
		suite.versionURL,
	)
	addressResponseJson, _ := json.Marshal(addressResponse)

	assert.Equal(suite.T(), http.StatusOK, respRetrieveAddressList.Code)
	assert.Contains(suite.T(), gjson.Get(bodyRetrieveAddressListJson, "_embedded.addresses").String(), string(addressResponseJson))

	// CINEMAS CONTEXT

	// Create Cinemas
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.POST("/addresses/:address_id/cinemas", handlers.CreateCinemas)

	cinemaCreate := requests.Cinema{
		UUID:        suite.cinemaUUID,
		Name:        "Sala Majestic IMAX 1",
		Description: "Sala IMAX com profundidade de audio",
		Capacity:    120,
	}

	cinemaCreateJson, _ := json.Marshal(cinemaCreate)
	addressUUIDCinemaRoute := fmt.Sprintf("/v1/addresses/%s/cinemas", suite.addressUUID.String())
	reqCinemasCreate, _ := http.NewRequest("POST", addressUUIDCinemaRoute, bytes.NewBuffer(cinemaCreateJson))
	respCinemasCreate := httptest.NewRecorder()

	suite.router.ServeHTTP(respCinemasCreate, reqCinemasCreate)

	assert.Equal(suite.T(), http.StatusCreated, respCinemasCreate.Code)
	assert.Equal(suite.T(), respCinemasCreate.Header().Get("Content-Type"), responses.HALHeaders["Content-Type"])

	// Retrieve Cinema
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/cinemas/:cinema_id", handlers.RetrieveCinema)

	cinemaUUIDRoute := fmt.Sprintf("/v1/cinemas/%v", suite.cinemaUUID.String())

	reqRetrieve, _ := http.NewRequest("GET", cinemaUUIDRoute, nil)
	respRetrieve := httptest.NewRecorder()
	suite.router.ServeHTTP(respRetrieve, reqRetrieve)

	bodyRetrieveCinemaJson := respRetrieve.Body.String()
	assert.Equal(suite.T(), http.StatusOK, respRetrieve.Code)
	assert.Equal(suite.T(), respRetrieve.Header().Get("Content-Type"), responses.JSONDefaultHeaders["Content-Type"])

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

	cinemaUpdateJson, _ := json.Marshal(cinemaUpdateRequest)
	reqCinemaUpdate, _ := http.NewRequest("PATCH", cinemaUUIDRoute, bytes.NewBuffer(cinemaUpdateJson))
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
	reqRetrieveCinemaList, _ := http.NewRequest("GET", addressCinemasListUUIDRoute, nil)
	respRetrieveCinemaList := httptest.NewRecorder()
	suite.router.ServeHTTP(respRetrieveCinemaList, reqRetrieveCinemaList)

	bodyRetrieveCinemaListJson := respRetrieveCinemaList.Body.String()

	addressCinemaListModel := models.Address{}
	if err := database.DB.Where(&models.Address{UUID: suite.addressUUID}).First(&addressCinemaListModel).Error; err != nil {
		fmt.Println("Address Not Found")
		return
	}
	cinemaModel, _ := models.NewCinema(
		suite.cinemaUUID,
		addressCinemaListModel.ID,
		cinemaCreate.Name,
		cinemaUpdateRequest.Description,
		cinemaUpdateRequest.Capacity,
	)

	cinemaResponse := responses.NewCinema(
		cinemaModel,
		addressResponse.Links.Self.HREF,
		suite.versionURL,
	)
	cinemaResponseJson, _ := json.Marshal(cinemaResponse)

	assert.Equal(suite.T(), http.StatusOK, respRetrieveCinemaList.Code)
	assert.Contains(suite.T(), gjson.Get(bodyRetrieveCinemaListJson, "_embedded.cinemas").String(), string(cinemaResponseJson))

	// MOVIES CONTEXT

	// Create Movies
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.POST("/movies", handlers.CreateMovies)

	ageRating := int64(14)
	published := true
	subtitled := false
	movieCreate := requests.Movie{
		UUID:        suite.movieUUID,
		Name:        "Back To The Recursion",
		Description: "Uma aventura no tempo usando técnicas avançadas de desenvolvimento de software",
		AgeRating:   &ageRating,
		Published:   &published,
		Subtitled:   &subtitled,
	}
	movieCreateJson, _ := json.Marshal(movieCreate)
	reqMoviesCreate, _ := http.NewRequest("POST", "/v1/movies", bytes.NewBuffer(movieCreateJson))
	respMoviesCreate := httptest.NewRecorder()
	suite.router.ServeHTTP(respMoviesCreate, reqMoviesCreate)

	assert.Equal(suite.T(), http.StatusCreated, respMoviesCreate.Code)

	// POSTERS CONTEXT

	// Upload Movie Poster
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.POST("/movies/:movie_id/posters", handlers.UploadMoviePoster)

	posterPath := "../../docs/assets/images/posters/back_to_the_recursion.png"
	posterFile, err := os.Open(posterPath)
	assert.NoError(suite.T(), err)
	defer posterFile.Close()

	fileInfo, _ := posterFile.Stat()
	fileBuffer := make([]byte, fileInfo.Size())
	posterFile.Read(fileBuffer)
	fileBytes := bytes.NewReader(fileBuffer)

	PosterRequestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(PosterRequestBody)
	posterFileHeader := make(textproto.MIMEHeader)
	posterFileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filepath.Base(posterPath)))
	posterFileHeader.Set("Content-Type", "image/png")
	posterFilePart, err := writer.CreatePart(posterFileHeader)
	assert.NoError(suite.T(), err)

	io.Copy(posterFilePart, fileBytes)

	fields := map[string]string{
		"uuid":            suite.posterUUID.String(),
		"name":            "Back To The Recursion",
		"alternativeText": "Uma aventura no tempo usando técnicas avançadas de desenvolvimento de software",
	}

	for key, value := range fields {
		err := writer.WriteField(key, value)
		assert.NoError(suite.T(), err)
	}
	writer.Close()

	uploadURL := fmt.Sprintf("/v1/movies/%s/posters", suite.movieUUID.String())
	reqUploadPoster, _ := http.NewRequest("POST", uploadURL, PosterRequestBody)
	reqUploadPoster.Header.Set("Content-Type", writer.FormDataContentType())
	respUploadPoster := httptest.NewRecorder()
	suite.router.ServeHTTP(respUploadPoster, reqUploadPoster)

	assert.Equal(suite.T(), http.StatusOK, respUploadPoster.Code)
}

func TestIntegrationSuccessfulSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuccesfulSuite))
}
