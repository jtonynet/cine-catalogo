package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/config"
)

func ConfigInject(cfg config.API) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("cfg", cfg)
		ctx.Next()
	}
}
