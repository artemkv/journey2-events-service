package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// port is a string in Gin format, e.g. ":8600"
func Serve(port string) {
	r := gin.Default()    // logging and recovery attached
	r.Use(cors.Default()) // allow all origins

	r.GET("/health", handleHealthCheck)

	r.Run(port)
}

// TODO: candidate to move to health package
func handleHealthCheck(c *gin.Context) {
	c.Status(http.StatusOK)
}
