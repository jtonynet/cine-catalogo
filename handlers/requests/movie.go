package requests

type Movie struct {
	Name        string `json:"name" binding:"required" hateoas:"placeholder:the name of movie;required:true"`
	Description string `json:"description" binding:"required" hateoas:"placeholder:the description of movie;required:true"`
	AgeRating   int64  `json:"age_rating" binding:"required" hateoas:"placeholder:the age rating of movie;required:true"`
	Subtitled   *bool  `json:"subtitled" binding:"required" hateoas:"placeholder:movie is subtitled;required:true"`
	Poster      string `json:"poster" binding:"required" hateoas:"placeholder:movie poster;required:true"`
}
