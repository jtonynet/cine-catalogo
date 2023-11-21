package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/config"
)

func Init(cfg config.API) {
	r := gin.Default()

	initializeRoutes(r, cfg)

	r.Run(fmt.Sprintf(":%s", cfg.Port))
}
