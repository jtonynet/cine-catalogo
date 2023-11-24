package responses

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/models"
)

type baseMovie struct {
	UUID        uuid.UUID `json:"uuid"`
	PosterUUID  uuid.UUID `json:"posterUUID"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AgeRating   int64     `json:"age_rating"`
	Subtitled   bool      `json:"subtitled"`
}

type Movie struct {
	baseMovie

	Links    HATEOASMovieLinks      `json:"_links,omitempty"`
	Embedded HATEOASMoviePosterList `json:"_embedded,omitempty"`

	Templates interface{} `json:"_templates"`
}

type MovieListItem struct {
	baseMovie

	Links HATEOASMovieLinks `json:"_links,omitempty"`
}

type HATEOASMovieLinks struct {
	Self              HATEOASLink `json:"self"`
	UpdateMovie       HATEOASLink `json:"update-movie"`
	UploadMoviePoster HATEOASLink `json:"upload-movie-poster"`
}

type HATEOASMoviePosterList struct {
	Posters []Poster `json:"posters,omitempty"`
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

		Links: HATEOASMovieLinks{
			Self:              HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s", versionURL, model.UUID.String())},
			UpdateMovie:       HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s", versionURL, model.UUID.String())},
			UploadMoviePoster: HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s/posters", versionURL, model.UUID.String())},
		},

		Templates: templates,
	}

	if len(model.Posters) > 0 {
		movie.Embedded = HATEOASMoviePosterList{
			Posters: []Poster{NewPoster(model.Posters[0], model.UUID, baseURL, versionURL, nil)},
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

		Links: HATEOASMovieLinks{
			Self:              HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s", versionURL, model.UUID.String())},
			UpdateMovie:       HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s", versionURL, model.UUID.String())},
			UploadMoviePoster: HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s/posters", versionURL, model.UUID.String())},
		},
	}

	if len(model.Posters) > 0 {
		movie.PosterUUID = model.Posters[0].UUID
	}

	return movie
}

type HATEOASMovieListLinks struct {
	Self         HATEOASLink `json:"self"`
	CreateMovies HATEOASLink `json:"create-movies"`
}

type HATEOASMovieAndPostersList struct {
	Movies  *[]MovieListItem `json:"movies"`
	Posters *[]Poster        `json:"posters"`
}

type MovieListResult struct {
	Embedded  HATEOASMovieAndPostersList `json:"_embedded"`
	Links     HATEOASMovieListLinks      `json:"_links"`
	Templates interface{}                `json:"_templates"`
}
