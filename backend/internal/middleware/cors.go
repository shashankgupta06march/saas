package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	config := cors.DefaultConfig()

	// Allow all origins - proper way without credentials conflict
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	config.ExposeHeaders = []string{"Content-Length", "Content-Type"}
	config.AllowCredentials = false // Must be false when allowing all origins
	config.MaxAge = 12 * 3600 // 12 hours

	return cors.New(config)
}
