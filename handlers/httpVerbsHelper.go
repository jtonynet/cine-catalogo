package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Option(ctx *gin.Context) {
	ctx.Header("Allow", "OPTIONS, GET, POST")
	ctx.Status(http.StatusNoContent)
}

func Head(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}
