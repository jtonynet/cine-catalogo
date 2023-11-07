package main

import (
	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/router"
)

func main() {

	cfg := config.Database{
		Host: "localhost",
		User: "api_user",
		Pass: "api_pass",
		DB:   "cine_catalog_db",
		Port: 5431,
	}
	database.Init(cfg)

	router.Init()
}
