package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/handlers/responses"
	"github.com/jtonynet/cine-catalogo/internal/hateoas"
	"github.com/tidwall/gjson"
)

func Option(ctx *gin.Context) {
	ctx.Header("Allow", "OPTIONS, GET, POST")
	ctx.Status(http.StatusNoContent)
}

func Head(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func templateFactory(
	rootURL string,
	templateParams []responses.HATEOASTemplateParams) (interface{}, error) {

	rootDocument := hateoas.NewRootDocument(rootURL)

	for _, param := range templateParams {
		resource, err := hateoas.NewResource(
			param.Name,
			param.ResourceURL,
			param.HTTPMethod,
		)
		if err != nil {
			// TODO: implements on future
			return nil, err
		}

		if param.RequestStruct != nil {
			resource.RequestToProperties(param.RequestStruct)
		}

		rootDocument.AddResource(resource)
	}

	rootEncoded, err := rootDocument.Encode()
	if err != nil {
		// TODO: implements on future
		return nil, err
	}

	templateString := gjson.Get(string(rootEncoded), "_templates").String()
	var templateJSON interface{}
	json.Unmarshal([]byte(templateString), &templateJSON)

	return templateJSON, nil
}
