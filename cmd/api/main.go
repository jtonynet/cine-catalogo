package main

import (
	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/router"
)

func main() {

	cfgDB := config.Database{
		Host: "localhost",
		User: "api_user",
		Pass: "api_pass",
		DB:   "cine_catalog_db",
		Port: 5432,
	}
	database.Init(cfgDB)

	cfgAPI := config.API{
		Name:       "catalogo",
		Port:       ":8080",
		TagVersion: "0.0.0",
		Env:        "dev",
		Host:       "http://localhost:8080/v1",
	}

	router.Init(cfgAPI)
}
