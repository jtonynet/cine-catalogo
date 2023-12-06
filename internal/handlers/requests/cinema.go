package requests

import "github.com/google/uuid"

type Cinema struct {
	UUID        uuid.UUID `json:"uuid" binding:"required" example:"292cb98c-62ab-49ef-8e23-dc793a86061d" hateoas:"name=uuid,placeholder=the uuid of movie,required=true,value=292cb98c-62ab-49ef-8e23-dc793a86061d,type=text"`
	Name        string    `json:"name" binding:"required" example:"Imax Majestic Room" hateoas:"placeholder=the name of cinema,required=true,value=Imax Majestic Room"`
	Description string    `json:"description" binding:"required" example:"Good holographic Imax 5D room" hateoas:"placeholder=the description of cinema,required=true,value=Good holographic Imax 5D room"`
	Capacity    int64     `json:"capacity" binding:"required" example:"120" hateoas:"placeholder=the country of capacity,required=true,value=120"`
}

type UpdateCinema struct {
	Name        string `json:"name" example:"5D Imax Majestic Room" hateoas:"placeholder=the name of cinema,required=true,value=Imax Majestic Room"`
	Description string `json:"description" example:"Majestic Very Good holographic Imax 5D room" hateoas:"placeholder=the description of cinema,required=true,value=Good holographic Imax 5D room"`
	Capacity    int64  `json:"capacity" example:"160" hateoas:"placeholder=the country of capacity,required=true,value=120"`
}
