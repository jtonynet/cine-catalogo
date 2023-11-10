package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/handlers/responses"
)

func RetrieveRootResources(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/hal+json")
	s := responses.Self{
		HRef: "http://localhost:8080/v1/",
	}
	l := responses.Links{Self: s}
	r := responses.RootResources{
		Links: l,
		ID:    "cine-catalogo",
		Name:  "Cine CataloGO",
	}

	ctx.JSON(http.StatusOK, r)

}
