package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	config := cors.DefaultConfig()

	// Allow all origins by using a custom function that always returns true
	// This is necessary because AllowCredentials = true prevents AllowAllOrigins = true
	config.AllowOriginFunc = func(origin string) bool {
		// Allow all origins (including null for iframes)
		return true
	}

	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	return cors.New(config)
}
