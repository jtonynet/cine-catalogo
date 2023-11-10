package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/handlers"
	"github.com/jtonynet/cine-catalogo/middlewares"
)

func initializeRoutes(r *gin.Engine) {
	basePath := "/v1"
	v1 := r.Group(basePath)
	v1.Use(middlewares.CORS())

	v1.GET("/", handlers.RetrieveRootResources)

	v1.GET("/addresses", handlers.RetrieveAddressList)
	v1.GET("/addresses/:addressId", handlers.RetrieveAddress)
	v1.POST("/addresses", handlers.CreateAddresses)
	v1.POST("addresses/:addressId/cinemas", handlers.CreateCinemas)

	v1.POST("/movies", handlers.CreateMovies)
	v1.GET("/movies", handlers.RetrieveMovieList)
}
