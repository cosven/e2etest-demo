package checker

import (
	"context"

	"github.com/pingcap/log"
)

type MetricsChecker struct {
	PrometheusUrl string
}

func (c *MetricsChecker) RunOnce(ctx context.Context) error {
	log.Info("Run once")
	return nil
}
