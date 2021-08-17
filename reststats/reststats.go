package reststats

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
}

func RequestCounter() gin.HandlerFunc {
	return func(c *gin.Context) {
		countRequest()
	}
}

type statsResult struct {
	RequestsTotal int `json:"requests_total"`
}

func HandleGetStats(c *gin.Context) {
	stats = getStats()

	result := &statsResult{
		RequestsTotal: stats.requestTotal,
	}

	c.JSON(http.StatusOK, result)
}

func CountRequestByEndpoint(endpoint string) {
	// TODO: implement
}

func UpdateResponseStats() {
	// TODO: implement
}
