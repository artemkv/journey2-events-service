package reststats

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CountRequestByEndpoint(endpoint string) {
	countRequestByEndpoint(endpoint)
}

func UpdateResponseStats(statusCode int) {
	updateResponseStats(statusCode)
}

func RequestCounter() gin.HandlerFunc {
	return func(c *gin.Context) {
		countRequest()
	}
}

type statsResult struct {
	RequestsTotal      int            `json:"requests_total"`
	RequestsByEndpoint map[string]int `json:"requests_by_endpoint"`
	ResponseStats      map[string]int `json:"responses_all"`
}

func HandleGetStats(c *gin.Context) {
	stats = getStats()

	result := &statsResult{
		RequestsTotal:      stats.requestTotal,
		RequestsByEndpoint: stats.requestsByEndpoint,
		ResponseStats:      stats.responseStats,
	}

	c.JSON(http.StatusOK, result)

	CountRequestByEndpoint("stats")
	UpdateResponseStats(c.Writer.Status())
}
