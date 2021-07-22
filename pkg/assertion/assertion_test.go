package assertion

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConsistentlyWithMetrics", func() {
	Context("should works fine", func() {
		const duration = 5 * time.Second

		It("within 1.1 times defaultConsistentlyWithMetricsDuration", func() {
			SetDefaultConsistentlyWithMetricsDuration(duration)
			SetDefaultConsistentlyWithMetricsPollingInterval(1 * time.Second)

			ConsistentlyWithMetrics(func() int { fmt.Println("checked"); return 1 }).Should(Equal(1))
		})
	})
})
