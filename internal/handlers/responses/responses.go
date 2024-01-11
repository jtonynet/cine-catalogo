package responses

import (
	"github.com/gin-gonic/gin"
)

type Header struct {
	key   string
	value string
}

func (h Header) Get(key string) string {
	if h.key == key {
		return h.value
	}
	return ""
}

var (
	JSONDefaultHeaders = []Header{{key: "Content-type", value: "application/json"}}
	HALHeaders         = []Header{{key: "Content-Type", value: "application/prs.hal-forms+json"}}
	MultipartFormData  = []Header{{key: "Content-Type", value: "multipart/form-data"}}
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

func setHeaders(ctx *gin.Context, headers []Header) {
	if headers == nil {
		headers = JSONDefaultHeaders
	}

	for _, header := range headers {
		ctx.Header(header.key, header.value)
	}
}
