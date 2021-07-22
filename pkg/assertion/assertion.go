package assertion

import (
	"time"
	_ "unsafe"

	"github.com/onsi/gomega"
)

var defaultConsistentlyWithMetricsDuration = time.Hour * 3
var defaultConsistentlyWithMetricsPollingInterval = time.Second * 30

func SetDefaultConsistentlyWithMetricsDuration(t time.Duration) {
	defaultConsistentlyWithMetricsDuration = t
}

func SetDefaultConsistentlyWithMetricsPollingInterval(t time.Duration) {
	defaultConsistentlyWithMetricsPollingInterval = t
}

func ConsistentlyWithMetrics(actual interface{}, intervals ...interface{}) gomega.AsyncAssertion {
	timeoutInterval := defaultConsistentlyWithMetricsDuration
	pollingInterval := defaultConsistentlyWithMetricsPollingInterval
	if len(intervals) > 0 {
		timeoutInterval = toDuration(intervals[0])
	}
	if len(intervals) > 1 {
		pollingInterval = toDuration(intervals[1])
	}

	return gomega.Consistently(actual, timeoutInterval, pollingInterval)
}

//go:linkname lookupStaticHost gomega.toDuration
func toDuration(input interface{}) time.Duration { return 42 * time.Second /* fake function body */ }
