package requests

import "mime/multipart"

type Poster struct {
	UUID            string `form:"uuid" binding:"required" hateoas:"name=uuid,placeholder=the uuid of poster,required=true,value=206dad85-cbcd-4b71-8fda-efd6ca87ebc7,type=text"`
	Name            string `form:"name" binding:"required" hateoas:"placeholder=the name of poster,required=true,value=Jaws,type=text"`
	AlternativeText string `form:"alternativeText" binding:"required" hateoas:"placeholder=the alternative text of poster,required=true,value=good poster,type=text"`

	// INFO: swaggerignore is a workaroud to fix swaggo bug. Swaggo dont recognize *multipart.FileHeader
	File *multipart.FileHeader `form:"file" format:"binary" binding:"required" swaggerignore:"true" hateoas:"placeholder=binary poster data,required=true,value=good poster,type=file"`
}
