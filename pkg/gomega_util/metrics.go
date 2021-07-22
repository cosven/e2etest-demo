package gomega_util

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/onsi/gomega/format"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func NewPrometheusClient(prometheusURL url.URL) (v1.API, error) {
	client, err := api.NewClient(api.Config{
		Address: prometheusURL.String(),
	})
	if err != nil {
		return nil, err
	}

	return v1.NewAPI(client), nil
}

type PromQLInstance struct {
	description string
	query       string
	api         *v1.API
}

// PromQL create a PromQLInstance that you could assert on.
// currently only PromQLEvaluatedToEmpty matcher is supported.
func PromQL(description string, query string, api *v1.API) PromQLInstance {
	return PromQLInstance{
		description: description,
		query:       query,
		api:         api,
	}
}

type promQLEvaluatedMatcher struct {
	evaluated   bool
	actual      model.Value
	description string
	promQL      string
}

// PromQLEvaluatedToEmpty evaluates a PromQL query against a cluster, and checks if it is:
//  - a Vector, and not empty
//  - a Scalar, and not empty
// if neither satisfied, the checker will fail. related checking logic is in promQLEvaluatedMatcher::checkPromQL.
func PromQLEvaluatedToEmpty() *promQLEvaluatedMatcher {
	return &promQLEvaluatedMatcher{}
}

func (matcher *promQLEvaluatedMatcher) Match(actual interface{}) (success bool, err error) {
	ql, ok := actual.(PromQLInstance)
	if !ok {
		return false, fmt.Errorf("PromQLEvaluatedToEmpty must be passed a util.PromQLInstance. Got:\n%s", format.Object(actual, 1))
	}

	result, value, err := matcher.checkPromQL(*ql.api, ql.query, time.Now())
	matcher.evaluated = err != nil
	matcher.actual = value
	matcher.promQL = ql.query
	matcher.description = ql.description
	return result, err
}

// checkPromQL copied from github.com/PingCAP-QE/metrics-checker with val returned
// checkPromQL checks if a query returns true, called by PromQLEvaluatedToEmpty.
func (matcher *promQLEvaluatedMatcher) checkPromQL(client v1.API, query string, ts time.Time) (result bool, value model.Value, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	val, _, err := client.Query(ctx, query, ts)
	if err != nil {
		return false, nil, err
	}
	// Ref: https://github.com/prometheus/prometheus/blob/76750d2a96df54226e85ac272d7ad5a547630240/rules/manager.go#L186-L206
	switch v := val.(type) {
	case model.Vector:
		return v.Len() > 0, val, nil
	case *model.Scalar:
		return v != nil, val, nil
	default:
		return false, val, errors.New("rule result is not a vector or scalar")
	}
}

func (matcher *promQLEvaluatedMatcher) FailureMessage(_ interface{}) (message string) {
	if !matcher.evaluated {
		return "PromQLInstance has not evaluated yet."
	} else {
		return format.Message(matcher.actual, fmt.Sprintf("to match the expected result of PromQL [%s] %s", matcher.description, matcher.promQL), "<non-empty Scalar or Vector>")
	}
}

func (matcher *promQLEvaluatedMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	if !matcher.evaluated {
		return "PromQLInstance has not evaluated yet."
	} else {
		return format.Message(matcher.actual, fmt.Sprintf("to match the expected result of PromQL [%s] %s", matcher.description, matcher.promQL), "<empty Scalar or Vector>")
	}
}
