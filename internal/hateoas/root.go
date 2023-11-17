package hateoas

import "github.com/pmoule/go2hal/halforms"

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

func NewRootDocument(href string) *root {
	return &root{
		document: halforms.NewDocument(href),
	}
}

func (r *root) AddResource(resource *resource) {
	r.resources = append(r.resources, *resource)

	r.document.AddTemplate(&resource.template)
	r.document.AddLink(resource.linkRelation)
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
