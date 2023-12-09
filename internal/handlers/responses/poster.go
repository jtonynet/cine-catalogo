package responses

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/models"
)

var MoviePosterContentType = "image/png"

type Poster struct {
	UUID            uuid.UUID `json:"uuid"`
	Name            string    `json:"name"`
	ContentType     string    `json:"contentType"`
	AlternativeText string    `json:"alternativeText"`
	Path            string    `json:"path"`

	Links HATEOASPosterLinks `json:"_links,omitempty"`

	Templates interface{} `json:"_templates,omitempty"`
}

type HATEOASPosterLinks struct {
	Self              HATEOASLink `json:"self"`
	Movie             HATEOASLink `json:"movie"`
	Image             HATEOASLink `json:"image"`
	UpdateMoviePoster HATEOASLink `json:"update-movie-poster"`
}

type PosterOption func(*Poster)

func NewPoster(
	model models.Poster,
	movieUUID uuid.UUID,
	movieLink,
	baseURL,
	versionURL string,
	options ...PosterOption,
) Poster {
	poster := Poster{

		UUID:            model.UUID,
		Name:            model.Name,
		ContentType:     model.ContentType,
		AlternativeText: model.AlternativeText,
		Path:            model.Path,

		Links: HATEOASPosterLinks{
			Self:              HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s/posters/%s", versionURL, movieUUID, model.UUID)},
			Movie:             HATEOASLink{HREF: movieLink},
			UpdateMoviePoster: HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s/posters/%s", versionURL, movieUUID, model.UUID)},
			Image:             HATEOASLink{HREF: fmt.Sprintf("%s/%s", baseURL, model.Path)},
		},
	}

	for _, opt := range options {
		opt(&poster)
	}

	return poster
}

func WithPosterTemplates(templates interface{}) PosterOption {
	return func(p *Poster) {
		p.Templates = templates
	}
}
