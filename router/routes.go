package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/handlers"
)

func initializeRoutes(r *gin.Engine) {
	basePath := "/v1"
	v1 := r.Group(basePath)

	v1.POST("/addresses", handlers.CreateAddresses)
	v1.GET("/addresses", handlers.RetrieveAddressList)
	v1.GET("/addresses/:uuid", handlers.RetrieveAddress)

	v1.POST("/cinemas", handlers.CreateCinema)

	v1.POST("/movies", handlers.CreateMovies)
	v1.GET("/movies", handlers.RetrieveMovieList)

}
