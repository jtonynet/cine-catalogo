package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/handlers"
)

func initializeRoutes(r *gin.Engine) {
	basePath := "/v1"
	v1 := r.Group(basePath)

	v1.GET("/", handlers.RootHandler)

	v1.POST("/addresses", handlers.CreateAddress)
	v1.GET("/addresses/:id", handlers.RetrieveAddress)

}
