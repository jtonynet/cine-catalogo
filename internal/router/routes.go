package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/config"
	docs "github.com/jtonynet/cine-catalogo/docs"
	"github.com/jtonynet/cine-catalogo/internal/handlers"
	"github.com/jtonynet/cine-catalogo/internal/middlewares"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initializeRoutes(r *gin.Engine, cfg config.API) {
	handlers.Init()

	basePath := "/v1"
	docs.SwaggerInfo.BasePath = basePath

	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL(fmt.Sprintf("%s%s", basePath, "/swagger/doc.json")),
		ginSwagger.DefaultModelsExpandDepth(-1))

	if cfg.MetricEnabled {
		handlers.ExposeMetrics(r, cfg)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Static("/web", cfg.StaticsDir)

	v1 := r.Group(basePath)

	v1.Use(middlewares.ConfigInject(cfg))
	v1.Use(middlewares.CORS())

	v1.GET("/", handlers.RetrieveRootResources)

	// Addresses
	v1.GET("/addresses", handlers.RetrieveAddressList)
	v1.GET("/addresses/:address_id", handlers.RetrieveAddress)

	v1.POST("/addresses", handlers.CreateAddresses)
	v1.OPTIONS("/addresses", handlers.Option)
	v1.HEAD("/addresses", handlers.Head)

	v1.PATCH("/addresses/:address_id", handlers.UpdateAddress)
	v1.DELETE("/addresses/:address_id", handlers.DeleteAddress)
	v1.OPTIONS("/addresses/:address_id", handlers.Option)
	v1.HEAD("/addresses/:address_id", handlers.Head)

	// Addresses Cinemas
	v1.GET("/addresses/:address_id/cinemas", handlers.RetrieveCinemaList)
	v1.POST("/addresses/:address_id/cinemas", handlers.CreateCinemas)
	v1.OPTIONS("/addresses/:address_id/cinemas", handlers.Option)
	v1.HEAD("/addresses/:address_id/cinemas", handlers.Head)

	// Cinemas
	v1.GET("/cinemas/:cinema_id", handlers.RetrieveCinema)
	v1.PATCH("/cinemas/:cinema_id", handlers.UpdateCinema)
	v1.DELETE("/cinemas/:cinema_id", handlers.DeleteCinema)
	v1.OPTIONS("/cinemas/:cinema_id", handlers.Option)
	v1.HEAD("/cinemas/:cinema_id", handlers.Head)

	// Movies
	v1.GET("/movies", handlers.RetrieveMovieList)
	v1.GET("/movies/:movie_id", handlers.RetrieveMovie)

	v1.POST("/movies", handlers.CreateMovies)
	v1.OPTIONS("/movies", handlers.Option)
	v1.HEAD("/movies", handlers.Head)

	v1.PATCH("/movies/:movie_id", handlers.UpdateMovie)
	v1.OPTIONS("/movies/:movie_id", handlers.Option)
	v1.HEAD("/movies/:movie_id", handlers.Head)

	// Movies Posters
	v1.GET("/movies/:movie_id/posters/:poster_id", handlers.RetrieveMoviePoster)
	v1.POST("/movies/:movie_id/posters", handlers.UploadMoviePoster)
	v1.OPTIONS("/movies/:movie_id/posters", handlers.Option)
	v1.HEAD("/movies/:movie_id/posters", handlers.Head)

	v1.PATCH("/movies/:movie_id/posters/:poster_id", handlers.UpdateMoviePoster)
	v1.OPTIONS("/movies/:movie_id/posters/:poster_id", handlers.Option)
	v1.HEAD("/movies/:movie_id/posters/:poster_id", handlers.Head)
}
