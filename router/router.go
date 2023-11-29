package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/config"
)

func InitMovies(cfg config.API) {
	r := gin.Default()

	initializeRoutes(r, cfg)

	r.Run(fmt.Sprintf(":%s", cfg.Port))
}
