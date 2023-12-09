package responses

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/internal/models"
)

type baseMovie struct {
	UUID        uuid.UUID `json:"uuid"`
	PosterLink  string    `json:"posterLink,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AgeRating   int64     `json:"age_rating"`
	Published   bool      `json:"published"`
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
	Poster            HATEOASLink `json:"poster"`
	UpdateMovie       HATEOASLink `json:"update-movie"`
	UploadMoviePoster HATEOASLink `json:"upload-movie-poster"`
}

type HATEOASMoviePosterList struct {
	Posters []Poster `json:"posters,omitempty"`
}

type MovieOption func(*Movie)

func NewMovie(
	model models.Movie,
	baseURL,
	versionURL string,
	options ...MovieOption,
) Movie {
	movieLink := fmt.Sprintf("%s/movies/%s", versionURL, model.UUID.String())
	movie := Movie{
		baseMovie: baseMovie{
			UUID:        model.UUID,
			Name:        model.Name,
			Description: model.Description,
			AgeRating:   model.AgeRating,
			Published:   model.Published,
			Subtitled:   model.Subtitled,
		},

		Links: HATEOASMovieLinks{
			Self:              HATEOASLink{HREF: movieLink},
			UpdateMovie:       HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s", versionURL, model.UUID.String())},
			UploadMoviePoster: HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s/posters", versionURL, model.UUID.String())},
		},
	}

	if len(model.Posters) > 0 {
		p := NewPoster(model.Posters[0], model.UUID, movieLink, baseURL, versionURL)
		movie.Embedded = HATEOASMoviePosterList{
			Posters: []Poster{p},
		}
		movie.Links.Poster = p.Links.Self
	}

	for _, opt := range options {
		opt(&movie)
	}

	return movie
}

func WithMovieTemplates(templates interface{}) MovieOption {
	return func(m *Movie) {
		m.Templates = templates
	}
}

func NewMovieListItem(
	model models.Movie,
	baseURL,
	versionURL string,
) MovieListItem {
	movieLink := fmt.Sprintf("%s/movies/%s", versionURL, model.UUID.String())
	movie := MovieListItem{
		baseMovie: baseMovie{
			UUID:        model.UUID,
			Name:        model.Name,
			Description: model.Description,
			AgeRating:   model.AgeRating,
			Published:   model.Published,
			Subtitled:   model.Subtitled,
		},

		Links: HATEOASMovieLinks{
			Self:              HATEOASLink{HREF: movieLink},
			UpdateMovie:       HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s", versionURL, model.UUID.String())},
			UploadMoviePoster: HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s/posters", versionURL, model.UUID.String())},
		},
	}

	if len(model.Posters) > 0 {
		p := NewPoster(model.Posters[0], model.UUID, movieLink, baseURL, versionURL)
		movie.Links.Poster = p.Links.Self
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
