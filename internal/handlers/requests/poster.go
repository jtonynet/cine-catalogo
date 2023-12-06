package requests

import "mime/multipart"

type Poster struct {
	UUID            string `form:"uuid" binding:"required" example:"2175d4e2-4d9c-411d-a986-08dc8f4e6a51" hateoas:"name=uuid,placeholder=the uuid of poster,required=true,value=2175d4e2-4d9c-411d-a986-08dc8f4e6a51,type=text"`
	Name            string `form:"name" binding:"required" example:"Spyder-Gopher" hateoas:"placeholder=the name of poster,required=true,value=Spyder-Gopher,type=text"`
	AlternativeText string `form:"alternativeText" binding:"required" example:"good poster of spyder gopher" hateoas:"placeholder=the alternative text of poster,required=true,value=good poster of spyder gopher,type=text"`

	// INFO: swaggerignore is a workaroud to fix swaggo bug. Swaggo dont recognize *multipart.FileHeader
	File *multipart.FileHeader `form:"file" format:"binary" binding:"required" swaggerignore:"true" hateoas:"placeholder=binary poster data,required=true,value=good poster,type=file"`
}

type UpdatePoster struct {
	Name            string `form:"name" example:"Spyder-Gopher" hateoas:"placeholder=the name of poster,required=true,value=Jaws,type=text"`
	AlternativeText string `form:"alternativeText" example:"good poster of spyder gopher AGAIN" hateoas:"placeholder=the alternative text of poster,required=true,value=good poster of spyder gopher AGAIN,type=text"`

	// INFO: swaggerignore is a workaroud to fix swaggo bug. Swaggo dont recognize *multipart.FileHeader
	File *multipart.FileHeader `form:"file" format:"binary" swaggerignore:"true" hateoas:"placeholder=binary poster data,required=true,value=good poster,type=file"`
}
