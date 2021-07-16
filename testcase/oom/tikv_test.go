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

	qe_metrics "github.com/PingCAP-QE/metrics-checker/pkg/metrics"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/pingcap/e2etest/pkg/ginkgo_helper"
	"github.com/pingcap/e2etest/pkg/workload"
	infra_resource "github.com/pingcap/test-infra/sdk/resource"
	_ "github.com/pingcap/test-infra/sdk/resource/impl/k8s"
)

var _ = ParameterizedGinkgoContainer("tikvoom", func(flagSet *flag.FlagSet) {
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

	Describe("Tikv OOM", func() {
		var (
			dbDSN         string
			prometheusURL *url.URL
		)
		ctx := suiteTestCtx
		BeforeEach(func() {
			r := ctx.Resource("tc")
			tc := r.(infra_resource.TiDBCluster)
			dbURL := Try(tc.ServiceURL(infra_resource.DBAddr)).(*url.URL)
			prometheusURL = Try(tc.ServiceURL(infra_resource.Prometheus)).(*url.URL)
			dbDSN = fmt.Sprintf("root@tcp(%s)/test", dbURL.Host)
			// dbDSN = "root@tcp(127.0.0.1:4000)/test"
			// prometheusUrl = "http://127.0.0.1:9090"
		})

		Context("Tikv under heavy workload", func() {
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

				checker := Try(qe_metrics.NewChecker(prometheusURL.String())).(*qe_metrics.Checker)
				rule := qe_metrics.Rule{
					Name:          "tikv should not restart.",
					PromQL:        `rate(process_start_time_seconds{component="tikv"}[1m]) != 0`,
					AlertCallback: func(rule qe_metrics.Rule) { Fail("One tikv restarted unexpected.") },
				}
				err := checker.RunRule(ctx2, rule, 10*time.Second)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

})
