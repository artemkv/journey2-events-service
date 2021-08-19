package reststats

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var version string = ""

func SetVersion(v string) {
	version = v
}

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
	Version              string         `json:"version"`
	Uptime               string         `json:"uptime"`
	TimeSinceLastRequest string         `json:"time_since_last_request"`
	RequestsTotal        int            `json:"requests_total"`
	RequestsByEndpoint   map[string]int `json:"requests_by_endpoint"`
	ResponseStats        map[string]int `json:"responses_all"`
}

func HandleGetStats(c *gin.Context) {
	stats = getStats()
	now := time.Now()

	result := &statsResult{
		Version:              version,
		Uptime:               getTimeDiffFormatted(stats.started, now),
		TimeSinceLastRequest: getTimeDiffFormatted(stats.previousRequestTime, now),
		RequestsTotal:        stats.requestTotal,
		RequestsByEndpoint:   stats.requestsByEndpoint,
		ResponseStats:        stats.responseStats,
	}

	c.JSON(http.StatusOK, result)

	CountRequestByEndpoint("stats")
	UpdateResponseStats(c.Writer.Status())
}

func getTimeDiffFormatted(start time.Time, end time.Time) string {
	return getTimeIntervalFormatted(end.Sub(start))
}

func getTimeIntervalFormatted(duration time.Duration) string {
	SECONDS_IN_DAY := 86400.0
	SECONDS_IN_HOUR := 3600.0
	SECONDS_IN_MINUTES := 60.0

	diff := duration.Seconds()

	days := math.Floor(diff / SECONDS_IN_DAY)
	diff = diff - days*SECONDS_IN_DAY

	hours := math.Floor(diff / SECONDS_IN_HOUR)
	diff = diff - hours*SECONDS_IN_HOUR

	minutes := math.Floor(diff / SECONDS_IN_MINUTES)
	diff = diff - minutes*SECONDS_IN_MINUTES

	seconds := math.Floor(diff)

	return fmt.Sprintf("%d.%d:%d:%d", int(days), int(hours), int(minutes), int(seconds))
}
