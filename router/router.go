package router

import (
	"github.com/gin-gonic/gin"
)

func Init() {
	r := gin.Default()

	initializeRoutes(r)

	r.Run(":8080")
}
