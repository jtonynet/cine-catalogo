package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/config"
)

func Init(cfg config.API) {
	r := gin.Default()

	initializeRoutes(r, cfg)

	r.Run(cfg.Port)
}
