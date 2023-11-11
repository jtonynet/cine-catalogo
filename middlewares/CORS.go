package middlewares

import "github.com/gin-gonic/gin"

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Enable access to HAL HATEOAS online tools for study and tests purpouses
		// TODO: Run HAL HATEOAS online tools on docker locally for more API security in near future
		// Access endpoint to test: https://hal-explorer.com/#theme=Dark&allHttpMethodsForLinks=true&uri=http://localhost:8080/v1/

		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://hal-explorer.com")

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
