package reststats

import "time"

var CURIOSITY = 1000
var CURIOSITY_FAILED = 100
var CURIOSITY_SLOW = 100
var SLOW_MS = 100
var QUICK_SEQUENCE_SIZE = 100

type statsData struct {
	started                  time.Time
	requestTotal             int
	requestsByEndpoint       map[string]int
	responseStats            map[string]int
	currentRequestTime       time.Time
	previousRequestTime      time.Time
	history                  []*responseStatsData
	historyOfFailed          []*responseStatsData
	historyOfSlow            []*responseStatsData
	shortestSequenceDuration time.Duration
}

type responseStatsData struct {
	time       time.Time
	url        string
	statusCode int
	duration   time.Duration
}

var stats = &statsData{
	started:                  time.Now(),
	requestTotal:             0,
	requestsByEndpoint:       map[string]int{},
	responseStats:            getEmptyCountsByStatusCodeMap(),
	currentRequestTime:       time.Now(),
	previousRequestTime:      time.Now(),
	history:                  make([]*responseStatsData, 0, CURIOSITY),
	historyOfFailed:          make([]*responseStatsData, 0, CURIOSITY_FAILED),
	historyOfSlow:            make([]*responseStatsData, 0, CURIOSITY_SLOW),
	shortestSequenceDuration: -1,
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

	stats.history = shiftAndPush(stats.history, responseStats, CURIOSITY)
	if statusCode >= 400 {
		stats.historyOfFailed = shiftAndPush(stats.historyOfFailed, responseStats, CURIOSITY_FAILED)
	}
	if duration >= time.Duration(SLOW_MS)*time.Millisecond {
		stats.historyOfSlow = shiftAndPush(stats.historyOfSlow, responseStats, CURIOSITY_SLOW)
	}

	updateCountsByStatusCodeMap(stats.responseStats, statusCode)

	if len(stats.history) >= QUICK_SEQUENCE_SIZE {
		lastSequenceDuration := stats.history[len(stats.history)-1].time.Sub(
			stats.history[len(stats.history)-QUICK_SEQUENCE_SIZE].time)
		if stats.shortestSequenceDuration == -1 || stats.shortestSequenceDuration > lastSequenceDuration {
			stats.shortestSequenceDuration = lastSequenceDuration
		}
	}
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

func shiftAndPush(slice []*responseStatsData, item *responseStatsData, maxLength int) []*responseStatsData {
	if len(slice) == maxLength {
		// TODO: study performance implications of this
		slice = slice[1:]
	}
	slice = append(slice, item)
	return slice
}
