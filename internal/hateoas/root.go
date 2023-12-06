package hateoas

import (
	"encoding/json"

	"github.com/pmoule/go2hal/halforms"
	"github.com/tidwall/gjson"
)

// WRAPPER FOR go2hal/hal AND go2hal/halforms TO SIMPLIFY USE
// https://rwcbook.github.io/hal-forms/#_the_hal_forms_media_type
// https://github.com/pmoule/go2hall
//
// HAL Client runs on docker image in port 4200
// http://localhost:4200/#uri=http://localhost:8080/v1/

type root struct {
	document  halforms.Document
	resources []resource
}

type TemplateParams struct {
	Name          string
	ResourceURL   string
	HTTPMethod    string
	ContentType   string
	RequestStruct interface{}
}

func NewRoot(href string) *root {
	return &root{
		document: halforms.NewDocument(href),
	}
}

func (r *root) AddResource(resource *resource) {
	r.resources = append(r.resources, *resource)

	r.document.AddTemplate(&resource.template)
	r.document.AddLink(resource.linkRelation)
}

func (r *root) GetDocument() halforms.Document {
	return r.document
}

func (r *root) Encode() ([]byte, error) {
	encoder := halforms.NewEncoder()
	bytes, err := encoder.ToJSON(r.document)
	if err != nil {
		// TODO Implements in future
		return nil, err
	}
	return bytes, nil
}

func TemplateFactory(
	baseURL string,
	TemplateParams []TemplateParams) (interface{}, error) {

	root := NewRoot(baseURL)

	for _, param := range TemplateParams {
		resource, err := NewResource(
			param.Name,
			param.ResourceURL,
			param.HTTPMethod,
			param.ContentType,
		)
		if err != nil {
			// TODO: implements on future
			return nil, err
		}

		if param.RequestStruct != nil {
			resource.RequestToProperties(param.RequestStruct)
		}

		root.AddResource(resource)
	}

	rootEncoded, err := root.Encode()
	if err != nil {
		// TODO: implements on future
		return nil, err
	}

	templateString := gjson.Get(string(rootEncoded), "_templates").String()
	var templateJSON interface{}
	json.Unmarshal([]byte(templateString), &templateJSON)

	return templateJSON, nil
}
