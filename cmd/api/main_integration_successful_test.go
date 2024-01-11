package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

	addressUUID uuid.UUID
	cinemaUUID  uuid.UUID
}

func (suite *IntegrationSuccesfulSuite) SetupSuite() {

	suite.cfg = setupConfig()
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)

	database.Init(suite.cfg.Database)

	suite.addressUUID, _ = uuid.Parse("9aa904a0-feed-4502-ace8-bf9dd0e23fb5")
	suite.cinemaUUID, _ = uuid.Parse("51276e29-940d-4d21-aa74-c0c4d3c5d632")
}

func (suite *IntegrationSuccesfulSuite) TearDownSuite() {
	query := fmt.Sprintf(`
	 DELETE FROM cinemas WHERE uuid in ('%v');
	 DELETE FROM addresses WHERE uuid in ('%v');
	 --DELETE FROM posters CASCADE;
	 --DELETE FROM movies CASCADE;`,
		suite.cinemaUUID.String(),
		suite.addressUUID.String())

	database.DB.Exec(query)

}

func setupConfig() *config.Config {
	cfg := config.Config{}

	cfg.API.Host = "catalogo-api-test"
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

func (suite *IntegrationSuccesfulSuite) SetupTest() {
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
}

func (suite *IntegrationSuccesfulSuite) TestV1CreateAddressesSuccessful() {
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
	reqCreate, _ := http.NewRequest("POST", "/v1/addresses", bytes.NewBuffer(addressCreateJson))
	respCreate := httptest.NewRecorder()
	suite.router.ServeHTTP(respCreate, reqCreate)
	assert.Equal(suite.T(), http.StatusCreated, respCreate.Code)

	// Retrieve Address
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/addresses/:address_id", handlers.RetrieveAddress)

	addressUUIDRoute := fmt.Sprintf("/v1/addresses/%s", suite.addressUUID.String())

	reqRetrieve, _ := http.NewRequest("GET", addressUUIDRoute, nil)
	respRetrieve := httptest.NewRecorder()
	suite.router.ServeHTTP(respRetrieve, reqRetrieve)

	bodyCreateJson := respRetrieve.Body.String()
	assert.Equal(suite.T(), http.StatusOK, respRetrieve.Code)
	assert.Equal(suite.T(), respRetrieve.Header().Get("Content-Type"), responses.JSONDefaultHeaders[0].Get("Content-type")) // Todo: Bad Smell, fix it

	assert.Equal(suite.T(), gjson.Get(bodyCreateJson, "uuid").String(), suite.addressUUID.String())
	assert.Equal(suite.T(), gjson.Get(bodyCreateJson, "country").String(), addressCreate.Country)
	assert.Equal(suite.T(), gjson.Get(bodyCreateJson, "state").String(), addressCreate.State)
	assert.Equal(suite.T(), gjson.Get(bodyCreateJson, "telephone").String(), addressCreate.Telephone)
	assert.Equal(suite.T(), gjson.Get(bodyCreateJson, "description").String(), addressCreate.Description)
	assert.Equal(suite.T(), gjson.Get(bodyCreateJson, "postalCode").String(), addressCreate.PostalCode)
	assert.Equal(suite.T(), gjson.Get(bodyCreateJson, "name").String(), addressCreate.Name)

	// Update Address
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.PATCH("/addresses/:address_id", handlers.UpdateAddress)

	addressUpdate := requests.UpdateAddress{
		Telephone: "1111-1111",
	}

	addressUpdateJson, _ := json.Marshal(addressUpdate)
	reqUpdate, _ := http.NewRequest("PATCH", addressUUIDRoute, bytes.NewBuffer(addressUpdateJson))
	respUpdate := httptest.NewRecorder()
	suite.router.ServeHTTP(respUpdate, reqUpdate)

	bodyUpdateJson := respUpdate.Body.String()
	assert.Equal(suite.T(), http.StatusOK, respUpdate.Code)
	assert.Equal(suite.T(), respUpdate.Header().Get("Content-Type"), responses.JSONDefaultHeaders[0].Get("Content-type")) // Todo: Bad Smell, fix it
	assert.Equal(suite.T(), gjson.Get(bodyUpdateJson, "telephone").String(), addressUpdate.Telephone)

	// Retrieve Address List
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
	suite.routesV1.GET("/addresses", handlers.RetrieveAddressList)

	reqRetrieveList, _ := http.NewRequest("GET", "/v1/addresses", nil)
	respRetrieveList := httptest.NewRecorder()
	suite.router.ServeHTTP(respRetrieveList, reqRetrieveList)

	bodyRetrieveListJson := respRetrieveList.Body.String()

	addressModel, _ := models.NewAddress(
		addressCreate.UUID,
		addressCreate.Country,
		addressCreate.State,
		addressUpdate.Telephone,
		addressCreate.Description,
		addressCreate.PostalCode,
		addressCreate.Name,
	)

	versionURL := fmt.Sprintf("%s/%s", suite.cfg.API.Host, "v1")
	addressResponse := responses.NewAddress(
		addressModel,
		versionURL,
	)
	addressResponseJson, _ := json.Marshal(addressResponse)

	assert.Equal(suite.T(), http.StatusOK, respRetrieveList.Code)
	assert.Contains(suite.T(), gjson.Get(bodyRetrieveListJson, "_embedded.addresses").String(), string(addressResponseJson))
}

func (suite *IntegrationSuccesfulSuite) TestV1CreateCinemasSuccessful() {
	// Create Cinemas
	suite.routesV1.POST("/addresses/:address_id/cinemas", handlers.CreateCinemas)

	cinema := requests.Cinema{
		UUID:        suite.cinemaUUID,
		Name:        "Sala Majestic IMAX 1",
		Description: "Sala IMAX com profundidade de audio",
		Capacity:    120,
	}

	cinemaJson, _ := json.Marshal(cinema)
	route := fmt.Sprintf("/v1/addresses/%s/cinemas", suite.addressUUID.String())
	reqCreate, _ := http.NewRequest("POST", route, bytes.NewBuffer(cinemaJson))
	respCreate := httptest.NewRecorder()

	context, _ := gin.CreateTestContext(respCreate)
	context.Request = reqCreate

	suite.router.ServeHTTP(respCreate, reqCreate)

	assert.Equal(suite.T(), http.StatusCreated, respCreate.Code)
	assert.Equal(suite.T(), respCreate.Header().Get("Content-Type"), responses.JSONDefaultHeaders[0].Get("Content-type")) // Todo: Bad Smell, fix it
}

func TestIntegrationSuccessfulSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuccesfulSuite))
}
