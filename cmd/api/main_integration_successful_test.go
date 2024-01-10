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

	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/internal/handlers"
	"github.com/jtonynet/cine-catalogo/internal/handlers/requests"
	"github.com/jtonynet/cine-catalogo/internal/middlewares"
)

type IntegrationSuccesfulSuite struct {
	suite.Suite
	cfg      *config.Config
	router   *gin.Engine
	routesV1 *gin.RouterGroup
}

func (suite *IntegrationSuccesfulSuite) SetupSuite() {

	suite.cfg = setupConfig()
	suite.router, suite.routesV1 = setupRouterAndGroup(suite.cfg.API)

	database.Init(suite.cfg.Database)
}

func (suite *IntegrationSuccesfulSuite) TearDownSuite() {
	fmt.Println(">>>End of IntegrationSuite!")
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

	gin.SetMode(gin.ReleaseMode)
	routes := gin.Default()

	routes.Use(middlewares.ConfigInject(cfg))

	return routes, routes.Group(basePath)
}

func (suite *IntegrationSuccesfulSuite) TestV1CreateAddressesSuccessful() {
	suite.routesV1.POST("/addresses", handlers.CreateAddresses)

	uniqueIDStr := "9aa904a0-feed-4502-ace8-bf9dd0e23fb5"
	uniqueID, _ := uuid.Parse(uniqueIDStr)

	address := requests.Address{
		UUID:        uniqueID,
		Country:     "BR",
		State:       "SP",
		Telephone:   "(11)0000-0000",
		Description: "Jardins Shoppings um dos mais belos de SP",
		PostalCode:  "1139050",
		Name:        "Jardins Shoppings",
	}

	addressJson, _ := json.Marshal(address)
	req, _ := http.NewRequest("POST", "/v1/addresses", bytes.NewBuffer(addressJson))
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuccesfulSuite))
}
