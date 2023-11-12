package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/handlers/requests"
	"github.com/jtonynet/cine-catalogo/internal/hateoas"
)

// HATEOAS flow controller - Is a good HATEOAS practice using HAL internelly
// Using hateoas import as wrapper to go2hal/hal and go2hal/halforms

func RetrieveRootResources(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/prs.hal-forms+json")

	rootURL := "http://localhost:8080/v1"
	root := hateoas.NewRoot(rootURL)

	createAddressesPost, err := hateoas.NewResource("createAddresses", "addresses", http.MethodPost)
	if err != nil {
		// TODO: implements on future
		return
	}
	createAddressesPost.RequestToProperties(requests.Address{})
	root.AddResource(createAddressesPost)

	createMoviesPost, err := hateoas.NewResource("createMovies", "movies", http.MethodPost)
	if err != nil {
		// TODO: implements on future
		return
	}
	createMoviesPost.RequestToProperties(requests.Movie{})
	root.AddResource(createMoviesPost)

	rootRendered, err := root.Render()
	if err != nil {
		// TODO: implements on future
		return
	}

	ctx.Data(http.StatusOK, "application/json", rootRendered)
}
