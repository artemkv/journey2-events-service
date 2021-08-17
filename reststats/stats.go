package reststats

type statsData struct {
	requestTotal       int
	requestsByEndpoint map[string]int
	responseStats      map[string]int
}

var stats = &statsData{
	requestTotal:       0,
	requestsByEndpoint: map[string]int{},
	responseStats: map[string]int{
		"1XX": 0,
		"2XX": 0,
		"3XX": 0,
		"4XX": 0,
		"5XX": 0,
	},
}

func getStats() *statsData {
	return stats
}

func countRequest() {
	stats.requestTotal++
}

func countRequestByEndpoint(endpoint string) {
	val, ok := stats.requestsByEndpoint[endpoint]
	if !ok {
		val = 0
	}
	stats.requestsByEndpoint[endpoint] = val + 1
}

func updateResponseStats(statusCode int) {
	if statusCode >= 500 {
		stats.responseStats["5XX"]++
	} else if statusCode >= 400 {
		stats.responseStats["4XX"]++
	} else if statusCode >= 300 {
		stats.responseStats["3XX"]++
	} else if statusCode >= 200 {
		stats.responseStats["2XX"]++
	} else {
		stats.responseStats["1XX"]++
	}
}
