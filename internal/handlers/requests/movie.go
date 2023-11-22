package requests

type Movie struct {
	Name        string `json:"name" binding:"required" hateoas:"placeholder=the name of movie,required=true,value=Gopher-Aranha,type=text"`
	Description string `json:"description" binding:"required" hateoas:"placeholder=the description of movie,required=true,value=O melhor super-herói de todos os tempos,type=text"`
	AgeRating   int64  `json:"age_rating" binding:"required" hateoas:"placeholder=the age rating of movie,required=true,value=12,type=text"`
	Subtitled   *bool  `json:"subtitled" binding:"required" hateoas:"placeholder=movie is subtitled,required=true,value=true,type=text"`
	Poster      string `json:"poster"`
}
