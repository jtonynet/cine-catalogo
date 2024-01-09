package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/internal/handlers"
	"github.com/jtonynet/cine-catalogo/internal/handlers/requests"
	"github.com/jtonynet/cine-catalogo/internal/middlewares"
)

func SetupEnvVars() (*config.Config, error) {
	// TODO: Mock config in future
	cfg, err := config.LoadConfig("../../.")
	if err != nil {
		return nil, err
	}

	return cfg, err
}

func SetupRouterV1() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	routes := gin.Default()
	cfg, err := SetupEnvVars()
	if err != nil {
		return nil
	}

	routes.Use(middlewares.ConfigInject(cfg.API))

	return routes
}

func TestV1CreateAddressesSucessful(t *testing.T) {
	cfg, _ := SetupEnvVars()
	cfg.Database.Host = "localhost"
	cfg.Database.MetricEnabled = false
	database.Init(cfg.Database)

	r := SetupRouterV1()
	r.POST("/v1/addresses", handlers.CreateAddresses)

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
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
