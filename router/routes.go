package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/handlers"
	"github.com/jtonynet/cine-catalogo/internal/middlewares"

	docs "github.com/jtonynet/cine-catalogo/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// INFO: To manage OPTION and HEAD verbs requests its necessary to implements HATEOAS HAL routes
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/OPTIONS#identifying_allowed_request_methods

func initializeRoutes(r *gin.Engine, cfg config.API) {
	r.Static("/web", cfg.StaticsDir)

	basePath := "/v1"

	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL(fmt.Sprintf("%s%s", basePath, "/swagger/doc.json")),
		ginSwagger.DefaultModelsExpandDepth(-1))

	docs.SwaggerInfo.BasePath = basePath
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group(basePath)

	v1.Use(middlewares.ConfigInject(cfg))
	v1.Use(middlewares.CORS())

	v1.GET("/", handlers.RetrieveRootResources)

	// Addresses
	v1.GET("/addresses", handlers.RetrieveAddressList)
	v1.GET("/addresses/:address_id", handlers.RetrieveAddress)
	v1.GET("/addresses/:address_id/cinemas", handlers.RetrieveCinemaList)

	v1.POST("/addresses", handlers.CreateAddresses)
	v1.OPTIONS("/addresses", handlers.Option)
	v1.HEAD("/addresses", handlers.Head)

	v1.POST("addresses/:address_id/cinemas", handlers.CreateCinemas)
	v1.OPTIONS("/addresses/:address_id/cinemas", handlers.Option)
	v1.HEAD("/addresses/:address_id/cinemas", handlers.Head)

	// Cinemas
	v1.GET("/cinemas/:cinema_id", handlers.RetrieveCinema)
	v1.POST("/movies", handlers.CreateMovies)

	// Movies
	v1.POST("/movies/:movie_id/posters", handlers.UploadMoviePoster)
	v1.OPTIONS("/movies", handlers.Option)
	v1.HEAD("/movies", handlers.Head)

	v1.GET("/movies", handlers.RetrieveMovieList)
	v1.GET("/movies/:movie_id", handlers.RetrieveMovie)
}

func initializeRoutesMovies(r *gin.Engine, cfg config.API) {
	r.Static("/web", cfg.StaticsDir)

	basePath := "/v1"

	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL(fmt.Sprintf("%s%s", basePath, "/swagger/doc.json")),
		ginSwagger.DefaultModelsExpandDepth(-1))

	docs.SwaggerInfo.BasePath = basePath
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group(basePath)

	v1.Use(middlewares.ConfigInject(cfg))
	v1.Use(middlewares.CORS())

	v1.GET("/", handlers.RetrieveRootResources)

	// Movies
	v1.POST("/movies/:movie_id/posters", handlers.UploadMoviePoster)
	v1.OPTIONS("/movies", handlers.Option)
	v1.HEAD("/movies", handlers.Head)

	v1.GET("/movies", handlers.RetrieveMovieList)
	v1.GET("/movies/:movie_id", handlers.RetrieveMovie)
}
