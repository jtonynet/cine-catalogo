package requests

import "github.com/google/uuid"

type Address struct {
	UUID        uuid.UUID `json:"uuid" binding:"required" hateoas:"name:uuidplaceholder=the uuid of address,required=true,value=2e61ddac-c3cc-46e9-ba88-0e86a790c924,type=text"`
	Country     string    `json:"country" binding:"required" hateoas:"placeholder=the country of address,required=true,value=BR"`
	State       string    `json:"state" binding:"required" hateoas:"placeholder=the state of address,required=true;value=RJ"`
	Telephone   string    `json:"telephone" binding:"required" hateoas:"placeholder=the telephone of address,required=true,value=9999-9999"`
	Description string    `json:"description" binding:"required" hateoas:"placeholder=the description of address,required=true,value=Bem localizado"`
	PostalCode  string    `json:"postalCode" binding:"required" hateoas:"placeholder=the postalCode of address,required=true,value=21940980"`
	Name        string    `json:"name" binding:"required" hateoas:"placeholder=the name of address,required=true,value=Botafogo Praia Center"`
}
