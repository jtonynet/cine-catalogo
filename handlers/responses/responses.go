package responses

import (
	"github.com/gin-gonic/gin"
)

type header struct {
	key   string
	value string
}

var (
	JSONDefaultHeaders = []header{{key: "Content-type", value: "application/json"}}
	HALHeaders         = []header{{key: "Content-Type", value: "application/prs.hal-forms+json"}}
)

func SendError(ctx *gin.Context, code int, msg string, headers []header) {
	setHeaders(ctx, headers)

	ctx.Header("Content-type", "application/json")
	ctx.JSON(code, gin.H{
		"message": msg,
		"code":    code,
	})
}

func SendSuccess(ctx *gin.Context, code int, op string, data interface{}, headers []header) {
	setHeaders(ctx, headers)

	ctx.JSON(code, data)
}

func setHeaders(ctx *gin.Context, headers []header) {
	if headers == nil {
		headers = JSONDefaultHeaders
	}

	for _, header := range headers {
		ctx.Header(header.key, header.value)
	}
}
