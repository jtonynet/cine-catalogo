package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/handlers"
	"github.com/jtonynet/cine-catalogo/middlewares"
)

// INFO: To manage OPTION and HEAD verbs requests its necessary to implements HATEOAS HAL routes
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/OPTIONS#identifying_allowed_request_methods

func initializeRoutes(r *gin.Engine, cfg config.API) {
	basePath := "/v1"
	v1 := r.Group(basePath)
	v1.Use(middlewares.ConfigInject(cfg))
	v1.Use(middlewares.CORS())

	v1.GET("/", handlers.RetrieveRootResources)

	v1.GET("/addresses", handlers.RetrieveAddressList)
	v1.GET("/addresses/:addressId", handlers.RetrieveAddress)

	v1.POST("/addresses", handlers.CreateAddresses)
	v1.OPTIONS("/addresses", handlers.Option)
	v1.HEAD("/addresses", handlers.Head)

	v1.POST("addresses/:addressId/cinemas", handlers.CreateCinemas)
	v1.OPTIONS("addresses/:addressId/cinemas", handlers.Option)
	v1.HEAD("addresses/:addressId/cinemas", handlers.Head)

	v1.POST("/movies", handlers.CreateMovies)
	v1.PUT("/movies/:movieId", handlers.UploadMoviePoster)
	v1.OPTIONS("/movies", handlers.Option)
	v1.HEAD("/movies", handlers.Head)

	v1.GET("/cinemas", handlers.RetrieveCinemaList)

	v1.GET("/movies", handlers.RetrieveMovieList)
}
