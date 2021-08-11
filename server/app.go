package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupRouter(router *gin.Engine) {
	router.Use(cors.Default()) // allow all origins

	router.GET("/health", handleHealthCheck)
}
