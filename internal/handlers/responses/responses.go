package responses

import (
	"github.com/gin-gonic/gin"
)

type Headers []Header

type Header struct {
	Key   string
	Value string
}

func (h Headers) Get(key string) string {
	for _, header := range h {
		if header.Key == key {
			return header.Value
		}
	}
	return ""
}

var (
	JSONDefaultHeaders = Headers{{Key: "Content-Type", Value: "application/json"}}
	HALHeaders         = Headers{{Key: "Content-Type", Value: "application/prs.hal-forms+json"}}
	MultipartFormData  = Headers{{Key: "Content-Type", Value: "multipart/form-data"}}
)

func SendError(ctx *gin.Context, code int, msg string, headers []Header) {
	setHeaders(ctx, headers)

	ctx.JSON(code, gin.H{
		"message": msg,
		"code":    code,
	})
}

func SendSuccess(ctx *gin.Context, code int, op string, data interface{}, headers []Header) {
	setHeaders(ctx, headers)

	ctx.JSON(code, data)
}

func setHeaders(ctx *gin.Context, headers Headers) {
	if headers == nil {
		headers = JSONDefaultHeaders
	}

	for _, header := range headers {
		ctx.Header(header.Key, header.Value)
	}
}
