package responses

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/models"
)

type Movie struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AgeRating   int64     `json:"age_rating"`
	Subtitled   bool      `json:"subtitled"`
	Poster      string    `json:"poster"`

	HATEOASListItemResult
}

type HATEOASMovieItemLinks struct {
	Self        HATEOASLink `json:"self"`
	UpdateMovie HATEOASLink `json:"update-movie"`
}

type HATEOASMovieListLinks struct {
	Self         HATEOASLink `json:"self"`
	CreateMovies HATEOASLink `json:"create-movies"`
}

type HATEOASMovieList struct {
	Movies *[]Movie `json:"movies"`
}

type MovieOption func(*Movie)

func NewMovie(
	model models.Movie,
	baseURL string,
	versionURL string,
	options ...MovieOption,
) *Movie {
	posterPath := fmt.Sprintf("%s/%s", baseURL, model.Poster)

	movie := &Movie{
		UUID:        model.UUID,
		Name:        model.Name,
		Description: model.Description,
		AgeRating:   model.AgeRating,
		Subtitled:   model.Subtitled,
		Poster:      posterPath,

		HATEOASListItemResult: HATEOASListItemResult{
			Links: HATEOASMovieItemLinks{
				Self: HATEOASLink{
					HREF: fmt.Sprintf("%s/movies/%s", baseURL, model.UUID.String()),
				},
				UpdateMovie: HATEOASLink{
					HREF: fmt.Sprintf("%s/movies/%s", baseURL, model.UUID.String()),
				},
			},
		},
	}

	for _, opt := range options {
		opt(movie)
	}

	return movie
}

func WithMovieTemplates(templates interface{}) MovieOption {
	return func(movie *Movie) {
		movie.Templates = templates
	}
}
