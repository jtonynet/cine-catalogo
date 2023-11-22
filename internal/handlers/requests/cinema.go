package requests

type Cinema struct {
	Name        string `json:"name" binding:"required" hateoas:"placeholder=the name of cinema,required=true,value=Sala Imax"`
	Description string `json:"description" binding:"required" hateoas:"placeholder=the description of cinema,required=true,value=Bom Cinema"`
	Capacity    int64  `json:"capacity" binding:"required" hateoas:"placeholder=the country of capacity,required=true,value=120"`
}
