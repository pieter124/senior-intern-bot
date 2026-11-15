package poller

import (
	"strings"
	"testing"
)

// BenchmarkPoll hits the live Greenhouse and Ashby APIs, so results depend on
// network conditions. Run with: go test ./poller -run '^$' -bench Poll -benchtime 3x
func BenchmarkPoll(b *testing.B) {
	if testing.Short() {
		b.Skip("hits live job board APIs")
	}
	p := New()
	for b.Loop() {
		postings, err := p.Poll()
		b.ReportMetric(float64(len(postings)), "postings")
		errCount := 0
		if err != nil {
			errCount = len(strings.Split(err.Error(), "\n"))
		}
		b.ReportMetric(float64(errCount), "errors")
	}
}