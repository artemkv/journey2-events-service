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

func HandleEndpointWithStats(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		handler(c)
		duration := time.Since(start)

		countRequestByEndpoint(c.Request.URL.Path)
		updateResponseStats(start, c.Request.RequestURI, c.Writer.Status(), duration)
	}
}

func HandleWithStats(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		handler(c)
		duration := time.Since(start)

		updateResponseStats(start, c.Request.RequestURI, c.Writer.Status(), duration)
	}
}

func UpdateResponseStatsOnRecover(start time.Time, url string, statusCode int) {
	updateResponseStats(start, url, statusCode, 0)
}

func RequestCounter() gin.HandlerFunc {
	return func(c *gin.Context) {
		countRequest()
	}
}

type statsResult struct {
	Version                             string              `json:"version"`
	Uptime                              string              `json:"uptime"`
	TimeSinceLastRequest                string              `json:"time_since_last_request"`
	RequestsTotal                       int                 `json:"requests_total"`
	RequestsByEndpoint                  map[string]int      `json:"requests_by_endpoint"`
	ShortestInterval100RequestsReceived string              `json:"shortest_interval_100_requests_received"`
	ResponsesAll                        map[string]int      `json:"responses_all"`
	ResponsesLast1000                   map[string]int      `json:"responses_last_1000"`
	RequestsLast10                      []*requestStatsData `json:"requests_last_10"`
	FailedRequestsLast10                []*requestStatsData `json:"failed_requests_last_10"`
	SlowRequestsLast10                  []*requestStatsData `json:"slow_requests_last_10"`
}

type requestStatsData struct {
	Url        string `json:"url"`
	StatusCode int    `json:"statusCode"`
	Duration   string `json:"duration"`
}

func HandleGetStats(c *gin.Context) {
	stats = getStats()
	now := time.Now()

	responsesHistory := getResponseHistory(stats.history)
	requestsLast10 := getLast10Requests(stats.history)
	failedRequestsLast10 := getLast10Requests(stats.historyOfFailed)
	slowRequestsLast10 := getLast10Requests(stats.historyOfSlow)

	result := &statsResult{
		Version:              version,
		Uptime:               getTimeDiffFormatted(stats.started, now),
		TimeSinceLastRequest: getTimeDiffFormatted(stats.previousRequestTime, now),
		RequestsTotal:        stats.requestTotal,
		RequestsByEndpoint:   stats.requestsByEndpoint,
		// TODO: last_1000_requests
		ShortestInterval100RequestsReceived: getTimeIntervalFormatted(stats.shortestSequenceDuration),
		ResponsesAll:                        stats.responseStats,
		ResponsesLast1000:                   responsesHistory,
		RequestsLast10:                      requestsLast10,
		FailedRequestsLast10:                failedRequestsLast10,
		SlowRequestsLast10:                  slowRequestsLast10,
	}

	c.JSON(http.StatusOK, result)
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

func getResponseHistory(history []*responseStatsData) map[string]int {
	responsesHistory := getEmptyCountsByStatusCodeMap()
	for _, v := range history {
		updateCountsByStatusCodeMap(responsesHistory, v.statusCode)
	}
	return responsesHistory
}

func getLast10Requests(history []*responseStatsData) []*requestStatsData {
	requestsLast10 := make([]*requestStatsData, 0, 10)
	if len(history) > 0 {
		idx := len(history) - 10
		if idx < 0 {
			idx = 0
		}
		for _, v := range history[idx:] {
			requestsLast10 = append(requestsLast10,
				&requestStatsData{
					Url:        v.url,
					StatusCode: v.statusCode,
					Duration:   getTimeIntervalFormatted(v.duration),
				})
		}
	}
	return requestsLast10
}
