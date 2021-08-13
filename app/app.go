package app

import (
	"artemkv.net/journey2/health"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {
	router.Use(cors.Default()) // allow all origins

	router.GET("/health", health.HandleHealthCheck)
	router.GET("/liveness", health.HandleLivenessCheck)
	router.GET("/readiness", health.HandleReadinessCheck)
}
