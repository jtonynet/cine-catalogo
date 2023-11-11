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
	ctx.Header("Content-Type", "application/prs.hal-forms+json")

	// Root defines:
	rootResourceUrl := "http://localhost:8080/v1"
	root := halforms.NewDocument(rootResourceUrl)

	//Address Resource defines:
	addressResourceURL := fmt.Sprintf("%s/%s", rootResourceUrl, "addresses")
	addressResourcePost, _ := halforms.NewHALFormsRelation("createAddresses", addressResourceURL) //skipped error handling
	addressResourcePostTemplate := halforms.NewTemplate()
	addressResourcePostTemplate.Method = http.MethodPost
	addressResourcePostTemplate.Target = addressResourceURL
	addressResourcePostTemplate.Key = "createAddresses"
	addressResourcePostTemplate.Title = ""

	// ...properties from template
	countryProp := halforms.NewProperty("country")
	countryProp.Prompt = "Country"
	countryProp.Placeholder = "the country of address"
	countryProp.Required = true
	addressResourcePostTemplate.Properties = append(
		addressResourcePostTemplate.Properties,
		countryProp,
	)

	stateProp := halforms.NewProperty("state")
	stateProp.Prompt = "State"
	stateProp.Placeholder = "the state address"
	stateProp.Required = true
	addressResourcePostTemplate.Properties = append(
		addressResourcePostTemplate.Properties,
		stateProp,
	)

	telProp := halforms.NewProperty("telephone")
	telProp.Prompt = "Telephone"
	telProp.Placeholder = "the telephone of address"
	telProp.Required = true
	addressResourcePostTemplate.Properties = append(
		addressResourcePostTemplate.Properties,
		telProp,
	)

	descProp := halforms.NewProperty("description")
	descProp.Prompt = "Description"
	descProp.Placeholder = "the description of address"
	descProp.Required = true
	addressResourcePostTemplate.Properties = append(
		addressResourcePostTemplate.Properties,
		descProp,
	)

	postalCodeProp := halforms.NewProperty("postalCode")
	postalCodeProp.Prompt = "PostalCode"
	postalCodeProp.Placeholder = "the postalCode of address"
	postalCodeProp.Required = true
	addressResourcePostTemplate.Properties = append(
		addressResourcePostTemplate.Properties,
		postalCodeProp,
	)

	nameProp := halforms.NewProperty("name")
	nameProp.Prompt = "Name"
	nameProp.Placeholder = "the name of address"
	nameProp.Required = true
	addressResourcePostTemplate.Properties = append(
		addressResourcePostTemplate.Properties,
		nameProp,
	)

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
