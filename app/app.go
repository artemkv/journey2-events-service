package app

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"artemkv.net/journey2/health"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {
	router.Use(requestLogger(log.StandardLogger()))
	router.Use(gin.CustomRecovery(recover))
	router.Use(cors.Default()) // allow all origins

	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	router.GET("/health", health.HandleHealthCheck)
	router.GET("/liveness", health.HandleLivenessCheck)
	router.GET("/readiness", health.HandleReadinessCheck)
	router.GET("/error", handleError)

	router.POST("/action", handleAction)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"err": "Not found"})
	})
}

func toSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func toBadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
}

func toInternalServerError(c *gin.Context, errText string) {
	c.JSON(http.StatusInternalServerError, gin.H{"err": errText})
}

func recover(c *gin.Context, err interface{}) {
	if errText, ok := err.(string); ok {
		toInternalServerError(c, errText)
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}

func requestLogger(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		message := fmt.Sprintf("%d %s %s",
			c.Writer.Status(),
			c.Request.Method,
			c.Request.URL.Path)

		logger.Info(message)
	}
}

func handleError(c *gin.Context) {
	panic("Test error")
}
