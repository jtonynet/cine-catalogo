package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmoule/go2hal/halforms"
)

// HAL Forms to HATEOAS flow controller - Is a good HATEOAS practice
// https://rwcbook.github.io/hal-forms/#_the_hal_forms_media_type
// https://github.com/pmoule/go2hall
// https://hal-explorer.com/#theme=Dark&allHttpMethodsForLinks=true&hkey0=Accept&hval0=application/prs.hal-forms+json&uri=http://localhost:8080/v1/

func RetrieveRootResources(ctx *gin.Context) {
	// Root defines:
	rootResourceUrl := "http://localhost:8080/v1"
	root := halforms.NewDocument(rootResourceUrl)

	//Address Resource defines:
	addressResourceURL := fmt.Sprintf("%s/%s", rootResourceUrl, "addresses")
	addressResourcePost, _ := halforms.NewHALFormsRelation("createAddresses", addressResourceURL) //skipped error handling
	addressResourcePostTemplate := halforms.NewTemplate()
	addressResourcePostTemplate.Method = http.MethodPost
	addressResourcePostTemplate.Target = addressResourceURL
	// ...properties from template
	root.AddTemplate(addressResourcePostTemplate)

	root.AddLink(addressResourcePost)

	encoder := halforms.NewEncoder()
	bytes, error := encoder.ToJSON(root)
	if error != nil {
		// TODO Implements in future
		return
	}

	ctx.Data(http.StatusOK, "application/json", bytes)
}
