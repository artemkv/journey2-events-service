package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleHealthCheck(c *gin.Context) {
	c.Status(http.StatusOK)
}
