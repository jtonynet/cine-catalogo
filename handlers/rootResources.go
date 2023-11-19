package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/handlers/requests"
	"github.com/jtonynet/cine-catalogo/internal/hateoas"
)

// HATEOAS flow controller - Is a good HATEOAS practice using HAL internelly
// Using hateoas import as wrapper to go2hal/hal and go2hal/halforms

func RetrieveRootResources(ctx *gin.Context) {
	cfg := ctx.MustGet("cfg").(config.API)

	ctx.Header("Content-Type", "application/prs.hal-forms+json")

	versionURL := fmt.Sprintf("%s/%s", cfg.Host, "v1")

	root := hateoas.NewRoot(versionURL)

	createAddressesPost, err := hateoas.NewResource(
		"create-addresses",
		fmt.Sprintf("%s/%s", versionURL, "addresses"),
		http.MethodPost,
	)
	if err != nil {
		// TODO: implements on future
		return
	}
	createAddressesPost.RequestToProperties(requests.Address{})
	root.AddResource(createAddressesPost)

	retrieveAddressListGet, err := hateoas.NewResource(
		"retrieve-address-list",
		fmt.Sprintf("%s/%s", versionURL, "addresses"),
		http.MethodGet,
	)
	if err != nil {
		// TODO: implements on future
		return
	}
	root.AddResource(retrieveAddressListGet)

	createMoviesPost, err := hateoas.NewResource(
		"create-movies",
		fmt.Sprintf("%s/%s", versionURL, "movies"),
		http.MethodPost,
	)
	if err != nil {
		// TODO: implements on future
		return
	}
	createMoviesPost.RequestToProperties(requests.Movie{})
	root.AddResource(createMoviesPost)

	retrieveMovieListGet, err := hateoas.NewResource(
		"retrieve-movie-list",
		fmt.Sprintf("%s/%s", versionURL, "movies"),
		http.MethodGet,
	)
	if err != nil {
		// TODO: implements on future
		return
	}
	root.AddResource(retrieveMovieListGet)

	rootJSON, err := root.Encode()
	if err != nil {
		// TODO: implements on future
		return
	}

	ctx.Data(http.StatusOK, "application/prs.hal-forms+json", rootJSON)
}
