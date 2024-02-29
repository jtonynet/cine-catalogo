package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	docs "github.com/jtonynet/cine-catalogo/api"
	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/handlers"
	"github.com/jtonynet/cine-catalogo/internal/middlewares"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initializeRoutes(
	cfg config.API,
	r *gin.Engine,
	addressHandler *handlers.AdrressHandler,
	cinemaHandler *handlers.CinemaHandler,
	movieHandler *handlers.MovieHandler,
	posterHandler *handlers.PosterHandler,
) {
	handlers.Init()

	basePath := "/v1"
	docs.SwaggerInfo.BasePath = basePath

	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL(fmt.Sprintf("%s%s", basePath, "/swagger/doc.json")),
		ginSwagger.DefaultModelsExpandDepth(-1))

	if cfg.MetricEnabled {
		initializeMetricsRoute(r, cfg)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Static("/web", cfg.StaticsDir)

	v1 := r.Group(basePath)

	v1.Use(middlewares.ConfigInject(cfg))
	v1.Use(middlewares.CORS())

	v1.GET("/", handlers.RetrieveRootResources)

	// Addresses
	v1.GET("/addresses", addressHandler.RetrieveAddressList)
	v1.GET("/addresses/:address_id", addressHandler.RetrieveAddress)

	v1.POST("/addresses", addressHandler.CreateAddresses)
	v1.OPTIONS("/addresses", handlers.Option)
	v1.HEAD("/addresses", handlers.Head)

	v1.PATCH("/addresses/:address_id", addressHandler.UpdateAddress)
	v1.DELETE("/addresses/:address_id", addressHandler.DeleteAddress)
	v1.OPTIONS("/addresses/:address_id", handlers.Option)
	v1.HEAD("/addresses/:address_id", handlers.Head)

	// Addresses Cinemas
	v1.GET("/addresses/:address_id/cinemas", cinemaHandler.RetrieveCinemaList)
	v1.POST("/addresses/:address_id/cinemas", cinemaHandler.CreateCinemas)
	v1.OPTIONS("/addresses/:address_id/cinemas", handlers.Option)
	v1.HEAD("/addresses/:address_id/cinemas", handlers.Head)

	// Cinemas
	v1.GET("/cinemas/:cinema_id", cinemaHandler.RetrieveCinema)
	v1.PATCH("/cinemas/:cinema_id", cinemaHandler.UpdateCinema)
	v1.DELETE("/cinemas/:cinema_id", cinemaHandler.DeleteCinema)
	v1.OPTIONS("/cinemas/:cinema_id", handlers.Option)
	v1.HEAD("/cinemas/:cinema_id", handlers.Head)

	// Movies
	v1.GET("/movies", movieHandler.RetrieveMovieList)
	v1.GET("/movies/:movie_id", movieHandler.RetrieveMovie)

	v1.POST("/movies", movieHandler.CreateMovies)
	v1.OPTIONS("/movies", handlers.Option)
	v1.HEAD("/movies", handlers.Head)

	v1.PATCH("/movies/:movie_id", movieHandler.UpdateMovie)
	v1.OPTIONS("/movies/:movie_id", handlers.Option)
	v1.HEAD("/movies/:movie_id", handlers.Head)

	// Movies Posters
	v1.GET("/movies/:movie_id/posters/:poster_id", posterHandler.RetrieveMoviePoster)
	v1.POST("/movies/:movie_id/posters", posterHandler.UploadMoviePoster)
	v1.OPTIONS("/movies/:movie_id/posters", handlers.Option)
	v1.HEAD("/movies/:movie_id/posters", handlers.Head)

	v1.PATCH("/movies/:movie_id/posters/:poster_id", posterHandler.UpdateMoviePoster)
	v1.OPTIONS("/movies/:movie_id/posters/:poster_id", handlers.Option)
	v1.HEAD("/movies/:movie_id/posters/:poster_id", handlers.Head)
}

func initializeMetricsRoute(r *gin.Engine, cfg config.API) {
	middlewares.InitPrometheus(r, cfg)

	r.Use(middlewares.Prometheus())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
