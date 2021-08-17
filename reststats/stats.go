package reststats

type statsData struct {
	requestTotal int
}

var stats = &statsData{
	requestTotal: 0,
}

func getStats() *statsData {
	return stats
}

func countRequest() {
	stats.requestTotal++
}
