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
	"github.com/jtonynet/cine-catalogo/internal/middlewares"
)

// Using gin engine or ginTestEngine for tests, context params bugfixes (move to ADR or another doc):
// https://medium.com/nerd-for-tech/testing-rest-api-in-go-with-testify-and-mockery-c31ea2cc88f9
// https://forum.golangbridge.org/t/how-to-test-gin-gonic-handler-function-within-a-function/33334/2
// https://github.com/gin-gonic/gin/issues/1292
// https://github.com/gin-gonic/gin/pull/2803
// https://github.com/gin-gonic/gin/issues/2778
// https://github.com/gin-gonic/gin/issues/2816

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
	database.DB.Exec(`
	DELETE FROM cinemas CASCADE;
	DELETE FROM addresses CASCADE;
	DELETE FROM posters CASCADE;
	DELETE FROM movies CASCADE;
  `)

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
	routes := gin.Default()

	routes.Use(middlewares.ConfigInject(cfg))

	return routes, routes.Group(basePath)
}

func (suite *IntegrationSuccesfulSuite) SetupTest() {
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)
}

func (suite *IntegrationSuccesfulSuite) TestV1CreateAddressesSuccessful() {
	suite.routesV1.POST("/addresses", handlers.CreateAddresses)

	address := requests.Address{
		UUID:        suite.addressUUID,
		Country:     "BR",
		State:       "SP",
		Telephone:   "(11)0000-0000",
		Description: "Jardins Shoppings um dos mais belos de SP",
		PostalCode:  "1139050",
		Name:        "Jardins Shoppings",
	}

	addressJson, _ := json.Marshal(address)
	reqCreate, _ := http.NewRequest("POST", "/v1/addresses", bytes.NewBuffer(addressJson))
	respCreate := httptest.NewRecorder()
	suite.router.ServeHTTP(respCreate, reqCreate)
	assert.Equal(suite.T(), http.StatusOK, respCreate.Code)

	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)

	suite.routesV1.GET("/addresses/:address_id", handlers.RetrieveAddress)

	retrieveAddressRoute := fmt.Sprintf("/v1/addresses/%s", suite.addressUUID.String())
	reqRetrieve, _ := http.NewRequest("GET", retrieveAddressRoute, nil)
	respRetrieve := httptest.NewRecorder()
	suite.router.ServeHTTP(respRetrieve, reqRetrieve)
	assert.Equal(suite.T(), http.StatusOK, respRetrieve.Code)

	bodyJson := respRetrieve.Body.String()
	assert.Equal(suite.T(), gjson.Get(bodyJson, "uuid").String(), suite.addressUUID.String())
	assert.Equal(suite.T(), gjson.Get(bodyJson, "country").String(), address.Country)
	assert.Equal(suite.T(), gjson.Get(bodyJson, "state").String(), address.State)
	assert.Equal(suite.T(), gjson.Get(bodyJson, "telephone").String(), address.Telephone)
	assert.Equal(suite.T(), gjson.Get(bodyJson, "description").String(), address.Description)
	assert.Equal(suite.T(), gjson.Get(bodyJson, "postalCode").String(), address.PostalCode)
	assert.Equal(suite.T(), gjson.Get(bodyJson, "name").String(), address.Name)
}

func (suite *IntegrationSuccesfulSuite) TestV1CreateCinemasSuccessful() {
	suite.routesV1.POST("/addresses/:address_id/cinemas", handlers.CreateCinemas)

	cinema := requests.Cinema{
		UUID:        suite.cinemaUUID,
		Name:        "Sala Majestic IMAX 1",
		Description: "Sala IMAX com profundidade de audio",
		Capacity:    120,
	}

	cinemaJson, _ := json.Marshal(cinema)
	route := fmt.Sprintf("/v1/addresses/%s/cinemas", suite.addressUUID.String())
	req, _ := http.NewRequest("POST", route, bytes.NewBuffer(cinemaJson))
	resp := httptest.NewRecorder()

	context, _ := gin.CreateTestContext(resp)
	context.Request = req

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func TestIntegrationSuccessfulSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuccesfulSuite))
}
