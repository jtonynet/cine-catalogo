package responses

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/models"
)

var MoviePosterContentType = "image/png"

type Movie struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AgeRating   int64     `json:"age_rating"`
	Subtitled   bool      `json:"subtitled"`

	HATEOASEmbeddedPosterItem
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
	Movies  *[]Movie              `json:"movies"`
	Posters *[]HATEOASPosterLinks `json:"posters"`
}

type HATEOASMovieItemEmbedded struct {
	Embedded HATEOASLink `json:"_embedded"`
}

type MovieOption func(*Movie)

func NewMovie(
	model models.Movie,
	baseURL string,
	versionURL string,
	options ...MovieOption,
) *Movie {

	movie := &Movie{
		UUID:        model.UUID,
		Name:        model.Name,
		Description: model.Description,
		AgeRating:   model.AgeRating,
		Subtitled:   model.Subtitled,

		HATEOASListItemResult: HATEOASListItemResult{
			Links: HATEOASMovieItemLinks{
				Self: HATEOASLink{
					HREF: fmt.Sprintf("%s/movies/%s", versionURL, model.UUID.String()),
				},
				UpdateMovie: HATEOASLink{
					HREF: fmt.Sprintf("%s/movies/%s", versionURL, model.UUID.String()),
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

type HATEOASEmbeddedPosterItem struct {
	Embedded *HATEOASPosterItem `json:"_embedded,omitempty"`
}

type HATEOASPosterItem struct {
	Poster HATEOASPosterItemLinks `json:"poster,omitempty"`
}

type HATEOASPosterItemLinks struct {
	Links HATEOASPosterLinks `json:"_links,omitempty"`
}

type HATEOASPosterLinks struct {
	HREF        string `json:"href,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

func WithMoviePosterEmbedded(baseURL, posterPath string) MovieOption {
	return func(movie *Movie) {
		movie.Embedded = NewPosterItem(baseURL, posterPath)
	}
}

func NewPosterItem(baseURL, posterPath string) *HATEOASPosterItem {
	return &HATEOASPosterItem{
		Poster: HATEOASPosterItemLinks{
			Links: *NewPosterLinks(baseURL, posterPath),
		},
	}
}

func NewPosterLinks(baseURL, posterPath string) *HATEOASPosterLinks {
	return &HATEOASPosterLinks{
		HREF:        fmt.Sprintf("%s/%s", baseURL, posterPath),
		ContentType: MoviePosterContentType,
	}
}
