package responses

import (
	"github.com/gin-gonic/gin"
)

var (
	JSONDefaultHeaders = map[string]string{"Content-Type": "application/json"}
	HALHeaders         = map[string]string{"Content-Type": "application/prs.hal-forms+json"}
	MultipartFormData  = map[string]string{"Content-Type": "multipart/form-data"}
)

func SendError(ctx *gin.Context, code int, msg string, headers map[string]string) {
	setHeaders(ctx, headers)

	ctx.JSON(code, gin.H{
		"message": msg,
		"code":    code,
	})
}

func SendSuccess(ctx *gin.Context, code int, op string, data interface{}, headers map[string]string) {
	setHeaders(ctx, headers)

	ctx.JSON(code, data)
}

func setHeaders(ctx *gin.Context, headers map[string]string) {
	if headers == nil {
		headers = JSONDefaultHeaders
	}

	for key, value := range headers {
		ctx.Header(key, value)
	}
}
