package requests

import "github.com/google/uuid"

type Cinema struct {
	UUID        uuid.UUID `json:"uuid" binding:"required" hateoas:"name=uuid,placeholder=the uuid of movie,required=true,value=292cb98c-62ab-49ef-8e23-dc793a86061d,type=text"`
	Name        string    `json:"name" binding:"required" hateoas:"placeholder=the name of cinema,required=true,value=Sala Imax"`
	Description string    `json:"description" binding:"required" hateoas:"placeholder=the description of cinema,required=true,value=Bom Cinema"`
	Capacity    int64     `json:"capacity" binding:"required" hateoas:"placeholder=the country of capacity,required=true,value=120"`
}
