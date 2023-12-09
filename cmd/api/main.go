package main

import (
	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/internal/router"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		//TODO: Implements in future
		return
	}

	database.Init(cfg.Database)
	router.Init(cfg.API)
}
