package oom_test

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pingcap/test-infra/sdk/resource"
	_ "github.com/pingcap/test-infra/sdk/resource/impl/k8s"

	assertion2 "github.com/pingcap/e2etest/pkg/assertion"
	"github.com/pingcap/e2etest/pkg/matcher"
	. "github.com/pingcap/e2etest/pkg/util"
	"github.com/pingcap/e2etest/pkg/workload"
)

var _ = ParameterizedTestCase("tikvoom", func(flagSet *flag.FlagSet) {
	var (
		duration      time.Duration
		threadsNUmber int
		tablesNumber  int
		paddingLength int
	)

	flagSet.DurationVar(&duration, "duration", time.Second*10, "Run duration.")
	flagSet.IntVar(&threadsNUmber, "threads", 1024, "Threads number.")
	flagSet.IntVar(&tablesNumber, "tables", 128, "Tables number.")
	flagSet.IntVar(&paddingLength, "padding", 4000000, "Row padding length.")

	Describe("TiKV OOM", func() {
		var (
			tc            resource.TiDBCluster
			dbDSN         string
			prometheusURL *url.URL
		)
		ctx := suiteTestCtx
		BeforeEach(func() {
			r := ctx.Resource("tc")
			tc = r.(resource.TiDBCluster)
			dbURL := Try(tc.ServiceURL(resource.DBAddr)).(*url.URL)
			prometheusURL = Try(tc.ServiceURL(resource.Prometheus)).(*url.URL)
			dbDSN = fmt.Sprintf("root@tcp(%s)/test", dbURL.Host)
			// dbDSN = "root@tcp(127.0.0.1:4000)/test"
			// prometheusUrl = "http://127.0.0.1:9090"
		})

		Context("TiKV under heavy workload", func() {
			It("should not oom", func() {
				db := Try(sql.Open("mysql", dbDSN)).(*sql.DB)
				ctx2, cancel := context.WithTimeout(ctx, duration)
				defer cancel()

				w := workload.AppendWorkload{
					DB:          db,
					Concurrency: threadsNUmber,
					Tables:      tablesNumber,
					PadLength:   paddingLength,
				}
				Try(w.Prepare())
				go w.Run(ctx2)

				prometheusURL, err := tc.ServiceURL(resource.Prometheus)
				Expect(err).ShouldNot(HaveOccurred())
				c, err := matcher.NewPrometheusClient(*prometheusURL)
				Expect(err).ShouldNot(HaveOccurred())

				assertion2.ConsistentlyWithMetrics(matcher.PromQL(
					"TiKV should not restart",
					`rate(process_start_time_seconds{component="tikv"}[1m]) != 0`, &c),
				).Should(matcher.PromQLEvaluatedToEmpty())
			})
		})
	})
})
