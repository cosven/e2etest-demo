module github.com/pingcap/e2etest

go 1.16

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/joho/godotenv v1.3.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.10.1
	github.com/pingcap/errors v0.11.4
	github.com/pingcap/log v0.0.0-20210625125904-98ed8e2eb1c7
	github.com/pingcap/test-infra/sdk v0.0.0-20210713041825-aad776bd52f9
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.29.0
	go.uber.org/zap v1.16.0
)

// replace github.com/PingCAP-QE/metrics-checker => ../../qe/metrics-checker/
