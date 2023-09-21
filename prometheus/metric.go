package prometheus

import (
	"context"
	"fmt"
	prom_v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"k8s.io/apimachinery/pkg/api/errors"
	"log"
	"strings"
	"time"
)

func fetchRange(ctx context.Context, api prom_v1.API, query string, bounds prom_v1.Range) Metric {
	fmt.Println(query)
	result, warnings, err := api.QueryRange(ctx, query, bounds)
	if len(warnings) > 0 {
		log.Println("fetchRange. Prometheus Warnings", strings.Join(warnings, ","))
	}
	if err != nil {
		return Metric{Err: err}
	}

	switch result.Type() {
	case model.ValMatrix:
		return Metric{Matrix: result.(model.Matrix)}
	}
	return Metric{Err: fmt.Errorf("invalid query, matrix expected: %s", query)}
}

func fetchRateRange(ctx context.Context, api prom_v1.API, metricName string, labels []string, grouping string, q *RangeQuery) Metric {
	var query string
	for i, labelsInstance := range labels {
		if i > 0 {
			query += " OR "
		}
		if grouping == "" {
			query += fmt.Sprintf("sum(%s(%s%s[%s]))", q.RateFunc, metricName, labelsInstance, q.RateInterval)
		} else {
			query += fmt.Sprintf("sum(%s(%s%s[%s])) by (%s)", q.RateFunc, metricName, labelsInstance, q.RateInterval, grouping)
		}
	}
	if len(labels) > 1 {
		query = fmt.Sprintf("(%s)", query)
	}
	return fetchRange(ctx, api, query, q.Range)
}

func fetchHistogramValues(ctx context.Context, api prom_v1.API, metricName, labels, grouping, rateInterval string, avg bool, quantiles []string, queryTime time.Time) (map[string]model.Vector, error) {
	queries := buildHistogramQueries(metricName, labels, grouping, rateInterval, avg, quantiles)
	histogram := make(map[string]model.Vector, len(queries))
	for k, query := range queries {
		log.Println("[Prom] fetchHistogramValues: ", query)
		result, warnings, err := api.Query(ctx, query, queryTime)
		if len(warnings) > 0 {
			log.Println("[Prom] fetchHistogramValues. Prometheus Warnings: ", strings.Join(warnings, ","))
		}
		if err != nil {
			return nil, errors.NewServiceUnavailable(err.Error())
		}
		histogram[k] = result.(model.Vector)
	}
	return histogram, nil
}

func buildHistogramQueries(metricName, labels, grouping, rateInterval string, avg bool, quantiles []string) map[string]string {
	queries := make(map[string]string)
	if avg {
		groupingAvg := ""
		if grouping != "" {
			groupingAvg = fmt.Sprintf(" by (%s)", grouping)
		}
		// sum(rate(my_histogram_sum{foo=bar}[5m])) by (baz) / sum(rate(my_histogram_count{foo=bar}[5m])) by (baz)
		query := fmt.Sprintf("sum(irate(%s_sum%s[%s]))%s / sum(irate(%s_count%s[%s]))%s",
			metricName, labels, rateInterval, groupingAvg, metricName, labels, rateInterval, groupingAvg)
		queries["avg"] = query
	}

	groupingQuantile := ""
	if grouping != "" {
		groupingQuantile = fmt.Sprintf(",%s", grouping)
	}
	for _, quantile := range quantiles {
		// histogram_quantile(0.5, sum(rate(my_histogram_bucket{foo=bar}[5m])) by (le,baz))
		query := fmt.Sprintf("histogram_quantile(%s, sum(irate(%s_bucket%s[%s])) by (le%s))",
			quantile, metricName, labels, rateInterval, groupingQuantile)
		queries[quantile] = query
	}

	return queries
}

func fetchHistogramRange(ctx context.Context, api prom_v1.API, metricName, labels, grouping string, q *RangeQuery) Histogram {
	queries := buildHistogramQueries(metricName, labels, grouping, q.RateInterval, q.Avg, q.Quantiles)
	histogram := make(Histogram, len(queries))
	for k, query := range queries {
		histogram[k] = fetchRange(ctx, api, query, q.Range)
	}
	return histogram
}
