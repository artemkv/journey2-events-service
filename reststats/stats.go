package reststats

import "time"

var CURIOSITY = 1000

type statsData struct {
	started             time.Time
	requestTotal        int
	requestsByEndpoint  map[string]int
	responseStats       map[string]int
	currentRequestTime  time.Time
	previousRequestTime time.Time
	history             []*responseStatsData
}

type responseStatsData struct {
	time       time.Time
	url        string
	statusCode int
	duration   time.Duration
}

var stats = &statsData{
	started:             time.Now(),
	requestTotal:        0,
	requestsByEndpoint:  map[string]int{},
	responseStats:       getEmptyCountsByStatusCodeMap(),
	currentRequestTime:  time.Now(),
	previousRequestTime: time.Now(),
	history:             make([]*responseStatsData, 0, CURIOSITY),
}

func getStats() *statsData {
	return stats
}

// TODO: not thread-safe!
func countRequest() {
	stats.requestTotal++
	stats.previousRequestTime = stats.currentRequestTime
	stats.currentRequestTime = time.Now()
}

// TODO: not thread-safe!
func countRequestByEndpoint(endpoint string) {
	val, ok := stats.requestsByEndpoint[endpoint]
	if !ok {
		val = 0
	}
	stats.requestsByEndpoint[endpoint] = val + 1
}

// TODO: not thread-safe!
func updateResponseStats(start time.Time, url string, statusCode int, duration time.Duration) {
	responseStats := &responseStatsData{
		time:       time.Now(),
		url:        url,
		statusCode: statusCode,
		duration:   duration,
	}

	if len(stats.history) == CURIOSITY {
		stats.history = stats.history[1:]
	}
	stats.history = append(stats.history, responseStats)

	updateCountsByStatusCodeMap(stats.responseStats, statusCode)
}

func getEmptyCountsByStatusCodeMap() map[string]int {
	return map[string]int{
		"1XX": 0,
		"2XX": 0,
		"3XX": 0,
		"4XX": 0,
		"5XX": 0,
	}
}

func updateCountsByStatusCodeMap(responseMap map[string]int, statusCode int) {
	if statusCode >= 500 {
		responseMap["5XX"]++
	} else if statusCode >= 400 {
		responseMap["4XX"]++
	} else if statusCode >= 300 {
		responseMap["3XX"]++
	} else if statusCode >= 200 {
		responseMap["2XX"]++
	} else {
		responseMap["1XX"]++
	}
}
