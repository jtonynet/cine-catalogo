package responses

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/models"
)

var MoviePosterContentType = "image/png"

type basePoster struct {
	UUID            uuid.UUID `json:"uuid"`
	MovieUUID       uuid.UUID `json:"movieUUID,omitempty"`
	Name            string    `json:"name"`
	ContentType     string    `json:"contentType"`
	AlternativeText string    `json:"alternativeText"`
	Path            string    `json:"path"`
}

type Poster struct {
	basePoster

	Links HATEOASPosterLinks `json:"_links,omitempty"`

	//Templates              interface{}         `json:"_templates,omitempty"`
}

type HATEOASPosterLinks struct {
	Self              HATEOASLink `json:"self"`
	Image             HATEOASLink `json:"image"`
	UpdateMoviePoster HATEOASLink `json:"update-movie-poster"`
	DeleteMoviePoster HATEOASLink `json:"delete-movie-poster"`
}

func NewPoster(
	model models.Poster,
	movieUUID uuid.UUID,
	baseURL,
	versionURL string,
	templates interface{},
) Poster {
	poster := Poster{
		basePoster: basePoster{
			UUID:            model.UUID,
			MovieUUID:       movieUUID,
			Name:            model.Name,
			ContentType:     model.ContentType,
			AlternativeText: model.AlternativeText,
			Path:            model.Path,
		},

		Links: NewPosterLinks(movieUUID, model.UUID, baseURL, versionURL, model.Path),

		//Templates: templates,
	}

	return poster
}

func NewPosterLinks(
	movieUUID,
	posterUUID uuid.UUID,
	baseURL,
	versionURL,
	posterPath string,
) HATEOASPosterLinks {
	return HATEOASPosterLinks{
		Self:              HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s/posters", baseURL, movieUUID)},
		UpdateMoviePoster: HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s/posters/%s", baseURL, movieUUID, posterUUID)},
		DeleteMoviePoster: HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s/posters/%s", baseURL, movieUUID, posterUUID)},
		Image:             HATEOASLink{HREF: fmt.Sprintf("%s/%s", baseURL, posterPath)},
	}
}
