package requests

type Movie struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	AgeRating   int64  `json:"age_rating" binding:"required"`
	Subtitled   *bool  `json:"subtitled" binding:"required"`
	Poster      string `json:"poster" binding:"required"`
}
