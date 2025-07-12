package middleware

import (
	"foodcook/internal/pkg/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	cfg := config.GetConfig()
	if cfg == nil {
		return cors.Default()
	}

	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
	}

	return cors.New(corsConfig)
}
