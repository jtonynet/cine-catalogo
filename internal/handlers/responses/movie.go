package responses

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/models"
)

var MoviePosterContentType = "image/png"

type BaseMovie struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AgeRating   int64     `json:"age_rating"`
	Subtitled   bool      `json:"subtitled"`
}

type Movie struct {
	BaseMovie

	Templates interface{} `json:"_templates,omitempty"`
	HATEOASEmbeddedPosterItem
}

type MovieListItem struct {
	BaseMovie

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
	versionURL string,
) Movie {
	movie := Movie{
		BaseMovie: BaseMovie{
			UUID:        model.UUID,
			Name:        model.Name,
			Description: model.Description,
			AgeRating:   model.AgeRating,
			Subtitled:   model.Subtitled,
		},

		HATEOASEmbeddedPosterItem: HATEOASEmbeddedPosterItem{
			Embedded: &HATEOASPosterItem{
				Poster: HATEOASPosterItemLinks{
					Links: HATEOASPosterLinks{
						HREF:        fmt.Sprintf("%s/%s", baseURL, model.Poster),
						ContentType: MoviePosterContentType,
					},
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
		BaseMovie: BaseMovie{
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
type HATEOASPosterItemLinks struct {
	Links HATEOASPosterLinks `json:"_links,omitempty"`
}

// swagger:ignore
type HATEOASPosterLinks struct {
	HREF        string `json:"href,omitempty"`
	ContentType string `json:"contentType,omitempty"`
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

// swagger:ignore
type MovieListResult struct {
	Embedded  HATEOASMovieAndPostersList `json:"_embedded"`
	Links     HATEOASMovieListLinks      `json:"_links"`
	Templates interface{}                `json:"_templates"`
}
