package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/handlers"
)

func Init(cfg config.API,
	addressHandler *handlers.AdrressHandler,
	cinemaHandler *handlers.CinemaHandler,
	movieHandler *handlers.MovieHandler,
	posterHandler *handlers.PosterHandler,
) {
	r := gin.Default()

	initializeRoutes(
		cfg,
		r,
		addressHandler,
		cinemaHandler,
		movieHandler,
		posterHandler,
	)

	r.Run(fmt.Sprintf(":%s", cfg.Port))
}
