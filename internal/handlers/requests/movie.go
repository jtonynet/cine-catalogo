package requests

import "github.com/google/uuid"

type Movie struct {
	UUID        uuid.UUID `json:"uuid" binding:"required" example:"206dad85-cbcd-4b71-8fda-efd6ca87ebc7" hateoas:"name=uuid,placeholder=the uuid of movie,required=true,value=206dad85-cbcd-4b71-8fda-efd6ca87ebc7,type=text"`
	Name        string    `json:"name" binding:"required" example:"Spyder-Gopher" hateoas:"placeholder=the name of movie,required=true,value=Spyder-Gopher,type=text"`
	Description string    `json:"description" binding:"required" example:"The best hero of all time" hateoas:"placeholder=the description of movie,required=true,value=The best hero of all time,type=text"`
	AgeRating   *int64    `json:"ageRating" binding:"required" example:"12" hateoas:"placeholder=the age rating of movie,required=true,value=12,type=text"`
	Published   *bool     `json:"published" binding:"required" example:"true" hateoas:"placeholder=movie is published,required=true,value=true,type=text"`
	Subtitled   *bool     `json:"subtitled" binding:"required" example:"true" hateoas:"placeholder=movie is subtitled,required=true,value=true,type=text"`
}

type UpdateMovie struct {
	Name        string `json:"name" example:"Spyder-Gopher" hateoas:"placeholder=the name of movie,required=true,value=Spyder-Gopher,type=text"`
	Description string `json:"description" example:"The best hero of all time AGAIN" hateoas:"placeholder=the description of movie,required=true,value=The best hero of all time AGAIN,type=text"`
	AgeRating   *int64 `json:"ageRating" example:"12" hateoas:"placeholder=the age rating of movie,required=true,value=12,type=text"`
	Published   *bool  `json:"published" example:"true" hateoas:"placeholder=movie is published,required=true,value=true,type=text"`
	Subtitled   *bool  `json:"subtitled" example:"true" hateoas:"placeholder=movie is subtitled,required=true,value=true,type=text"`
}
