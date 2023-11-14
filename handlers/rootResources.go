package handlers

import (
	"fmt"
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

	// ---------
	// TODO:
	// Implements WithRequest option
	// hateoas.NewResource(
	// 	"createAddresses",
	// 	"addresses",
	// 	http.MethodPost,
	// 	hateoas.WithRequest(requests.Address{}),
	// )
	//---------
	createAddressesPost, err := hateoas.NewResource(
		"createAddresses",
		fmt.Sprintf("%s/%s", rootURL, "addresses"),
		http.MethodPost,
	)
	if err != nil {
		// TODO: implements on future
		return
	}
	createAddressesPost.RequestToProperties(requests.Address{})
	root.AddResource(createAddressesPost)

	retrieveAddressListGet, err := hateoas.NewResource(
		"retrieveAddresses",
		fmt.Sprintf("%s/%s", rootURL, "addresses"),
		http.MethodGet,
	)
	if err != nil {
		// TODO: implements on future
		return
	}
	root.AddResource(retrieveAddressListGet)

	createMoviesPost, err := hateoas.NewResource(
		"createMovies",
		fmt.Sprintf("%s/%s", rootURL, "movies"),
		http.MethodPost,
	)
	if err != nil {
		// TODO: implements on future
		return
	}
	createMoviesPost.RequestToProperties(requests.Movie{})
	root.AddResource(createMoviesPost)

	retrieveMovieListGet, err := hateoas.NewResource(
		"retrieveMovieList",
		fmt.Sprintf("%s/%s", rootURL, "movies"),
		http.MethodGet,
	)
	if err != nil {
		// TODO: implements on future
		return
	}
	root.AddResource(retrieveMovieListGet)

	rootJSON, err := root.ToJSON()
	if err != nil {
		// TODO: implements on future
		return
	}

	ctx.Data(http.StatusOK, "application/json", rootJSON)
}
