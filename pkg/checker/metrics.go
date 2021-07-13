package checker

import (
	"context"

	"github.com/ngaut/log"
)

type MetricsChecker struct {
	PrometheusUrl string
}

func (c *MetricsChecker) RunOnce(ctx context.Context) error {
	log.Info("Run once")
	return nil
}
