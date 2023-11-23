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
	HATEOASMovieEmbeddedPosterItem
}

type MovieListItem struct {
	baseMovie

	MovieListItemResult
}

type HATEOASMovieItemLinks struct {
	Self              HATEOASLink `json:"self"`
	UpdateMovie       HATEOASLink `json:"update-movie"`
	UploadMoviePoster HATEOASLink `json:"upload-movie-poster"`
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

type MovieListItemResult struct {
	Links HATEOASMovieItemLinks `json:"_links"`
}

type HATEOASMovieEmbeddedPosterItem struct {
	Embedded *HATEOASMoviePosterItem `json:"_embedded,omitempty"`
}

type HATEOASMoviePosterItem struct {
	Poster Poster `json:"poster,omitempty"`
}

func NewMovie(
	model models.Movie,
	baseURL,
	versionURL string,
	templates interface{},
) Movie {
	movie := Movie{
		baseMovie: baseMovie{
			UUID:        model.UUID,
			Name:        model.Name,
			Description: model.Description,
			AgeRating:   model.AgeRating,
			Subtitled:   model.Subtitled,
		},

		Templates: templates,
	}

	if len(model.Posters) > 0 {
		movie.Embedded = &HATEOASMoviePosterItem{
			Poster: NewPoster(model.Posters[0], model.UUID, baseURL, versionURL, nil),
		}
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
				UploadMoviePoster: HATEOASLink{
					HREF: fmt.Sprintf("%s/movies/%s/posters", versionURL, model.UUID.String()),
				},
			},
		},
	}

	return movie
}

type MovieListResult struct {
	Embedded  HATEOASMovieAndPostersList `json:"_embedded"`
	Links     HATEOASMovieListLinks      `json:"_links"`
	Templates interface{}                `json:"_templates"`
}
