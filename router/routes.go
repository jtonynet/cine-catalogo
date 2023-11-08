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
	v1.GET("/addresses", handlers.RetrieveAddressList)
	v1.GET("/addresses/:uuid", handlers.RetrieveAddress)

	v1.POST("/movies", handlers.CreateMovie)
	v1.GET("/movies", handlers.RetrieveMovieList)

}
