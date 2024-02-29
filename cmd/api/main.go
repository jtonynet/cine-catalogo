package main

import (
	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/database"
	"github.com/jtonynet/cine-catalogo/internal/handlers"
	"github.com/jtonynet/cine-catalogo/internal/router"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic("cannot load environment variables")
	}

	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		panic("cannot connect to database")
	}

	addressHandler := handlers.NewAddressHandler(db)
	cinemaHandler := handlers.NewCinemaHandler(db)
	movieHandler := handlers.NewMovieHandler(db)
	posterHandler := handlers.NewPosterHandler(db)

	router.Init(
		cfg.API,
		addressHandler,
		cinemaHandler,
		movieHandler,
		posterHandler,
	)
}
