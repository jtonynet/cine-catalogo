package responses

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/models"
)

type baseMovie struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AgeRating   int64     `json:"age_rating"`
	Subtitled   bool      `json:"subtitled"`
}

type Movie struct {
	baseMovie

	Templates interface{} `json:"_templates,omitempty"`
	HATEOASEmbeddedPosterItem
}

type MovieListItem struct {
	baseMovie

	MovieListItemResult
}

type HATEOASMovieItemLinks struct {
	Self        HATEOASLink `json:"self"`
	UpdateMovie HATEOASLink `json:"update-movie"`
}

type HATEOASMovieListLinks struct {
	Self         HATEOASLink `json:"self"`
	CreateMovies HATEOASLink `json:"create-movies"`
}

type HATEOASMovieAndPostersList struct {
	Movies  *[]MovieListItem      `json:"movies"`
	Posters *[]HATEOASPosterLinks `json:"posters"`
}

type HATEOASMovieItemEmbedded struct {
	Embedded HATEOASLink `json:"_embedded"`
}

type MovieOption func(*Movie)

type MovieListItemResult struct {
	Links HATEOASMovieItemLinks `json:"_links"`
}

func NewMovie(
	model models.Movie,
	templates interface{},
	baseURL string,
) Movie {
	movie := Movie{
		baseMovie: baseMovie{
			UUID:        model.UUID,
			Name:        model.Name,
			Description: model.Description,
			AgeRating:   model.AgeRating,
			Subtitled:   model.Subtitled,
		},

		// TODO: CHANGE TO EMBEDDED POSTER ENTITY
		HATEOASEmbeddedPosterItem: HATEOASEmbeddedPosterItem{
			Embedded: &HATEOASPosterItem{
				// TODO:
				//Poster: NewPoster( model models.Poster, baseURL string, posterPath string, templates interface{})

				Poster: HATEOASPosterItemLinks{
					Links: NewPosterLinks(model.UUID, uuid.New(), baseURL, model.Poster),
				},
			},
		},

		Templates: templates,
	}

	return movie
}

func NewMovieListItem(
	model models.Movie,
	baseURL,
	versionURL string,
) MovieListItem {
	movie := MovieListItem{
		baseMovie: baseMovie{
			UUID:        model.UUID,
			Name:        model.Name,
			Description: model.Description,
			AgeRating:   model.AgeRating,
			Subtitled:   model.Subtitled,
		},

		MovieListItemResult: MovieListItemResult{
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

	return movie
}

type HATEOASEmbeddedPosterItem struct {
	Embedded *HATEOASPosterItem `json:"_embedded,omitempty"`
}

type HATEOASPosterItem struct {
	Poster HATEOASPosterItemLinks `json:"poster,omitempty"`
}

// swagger:ignore

// swagger:ignore
type MovieListResult struct {
	Embedded  HATEOASMovieAndPostersList `json:"_embedded"`
	Links     HATEOASMovieListLinks      `json:"_links"`
	Templates interface{}                `json:"_templates"`
}
