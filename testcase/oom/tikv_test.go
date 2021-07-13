package oom_test

import (
	"context"
	"database/sql"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pingcap/e2etest/pkg/checker"
	"github.com/pingcap/e2etest/pkg/workload"
	infra_sdk "github.com/pingcap/test-infra/sdk/core"
)

var _ = Describe("Tikv OOM", func() {
	var (
		dbUrl         string
		prometheusUrl string
		ctx           infra_sdk.TestContext
		duration      time.Duration
	)

	BeforeEach(func() {
		dbUrl = "root@tcp(127.0.0.1:4000)/test"
		prometheusUrl = "http://127.0.0.1:9090"
		duration = time.Second * 10
		err := error(nil)
		ctx, err = infra_sdk.BuildContext()
		Expect(err).NotTo(HaveOccurred())
	})

	Context("Tikv under heavy workload", func() {
		It("should not oom", func() {
			db, err := sql.Open("mysql", dbUrl)
			Expect(err).NotTo(HaveOccurred())

			ctx, cancel := context.WithTimeout(ctx, duration)
			defer cancel()

			w := workload.AppendWorkload{
				DB:          db,
				Concurrency: 1,
				Tables:      1,
				PadLength:   4000000,
			}
			err = w.Prepare()
			Expect(err).NotTo(HaveOccurred())
			go w.Run(ctx)

			c := checker.MetricsChecker{PrometheusUrl: prometheusUrl}
			Consistently(func() error {
				return c.RunOnce(ctx)
			}, duration, time.Second*1).ShouldNot(HaveOccurred())
		})
	})
})
