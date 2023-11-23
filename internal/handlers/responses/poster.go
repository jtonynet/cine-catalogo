package responses

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/cine-catalogo/models"
)

var MoviePosterContentType = "image/png"

type basePoster struct {
	UUID            uuid.UUID `gorm:"type:uuid;unique;not null"`
	Name            string
	ContentType     string
	AlternativeText string
	Path            string
}

type Poster struct {
	basePoster

	Templates interface{} `json:"_templates,omitempty"`
}

type HATEOASPosterItemLinks struct {
	Links *HATEOASPosterLinks `json:"_links,omitempty"`
}

type HATEOASPosterLinks struct {
	Self         HATEOASLink `json:"self"`
	Image        HATEOASLink `json:"image"`
	UpdatePoster HATEOASLink `json:"update-poster"`
	DeletePoster HATEOASLink `json:"delete-poster"`
}

func NewPoster(
	model models.Poster,
	movieUUID uuid.UUID,
	baseURL,
	versionURL string,
	templates interface{},
) Poster {
	poster := Poster{
		basePoster{
			UUID:            model.UUID,
			Name:            model.Name,
			ContentType:     model.ContentType,
			AlternativeText: model.AlternativeText,
			Path:            model.Path,
		},

		NewPosterLinks(movieUUID, model.UUID, baseURL, versionURL, model.Path),

		//templates,
	}

	return poster
}

func NewPosterLinks(
	movieUUID,
	posterUUID uuid.UUID,
	baseURL,
	versionURL,
	posterPath string,
) *HATEOASPosterLinks {
	return &HATEOASPosterLinks{
		Self:         HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s/posters", baseURL, movieUUID)},
		UpdatePoster: HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s/posters/%s", baseURL, movieUUID, posterUUID)},
		DeletePoster: HATEOASLink{HREF: fmt.Sprintf("%s/movies/%s/posters/%s", baseURL, movieUUID, posterUUID)},
		Image:        HATEOASLink{HREF: fmt.Sprintf("%s/%s", baseURL, posterPath)},
	}
}
