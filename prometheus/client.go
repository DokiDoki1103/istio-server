package prometheus

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	prom_v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"istio-server/config"
	"istio-server/graph"
	"k8s.io/apimachinery/pkg/api/errors"
	"os"
	"time"
)

var (
	promClient *Client
)

func GetPrometheusClient() *Client {

	return promClient
}

type Client struct {
	PromClientInterface
	p8s api.Client
	api prom_v1.API
	ctx context.Context
}

func NewClient() (*Client, error) {
	var url = os.Getenv("PROM_URL")
	if url == "" {
		url = "http://prometheus.istio-system.svc.cluster.local:9090"
	}
	cfg := config.PrometheusConfig{
		URL: url,
	}
	return NewClientForConfig(cfg)
}

func NewClientForConfig(cfg config.PrometheusConfig) (*Client, error) {
	clientConfig := api.Config{Address: cfg.URL}
	p8s, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, errors.NewServiceUnavailable(err.Error())
	}
	client := &Client{p8s: p8s, api: prom_v1.NewAPI(p8s), ctx: context.Background()}
	promClient = client
	return client, nil
}

type PromClientInterface interface {
	GenGraph(ns string, app string) (g graph.Graph, err error)

	FetchHistogramRange(metricName, labels, grouping string, q *RangeQuery) Histogram
	FetchHistogramValues(metricName, labels, grouping, rateInterval string, avg bool, quantiles []string, queryTime time.Time) (map[string]model.Vector, error)

	FetchRange(metricName, labels, grouping, aggregator string, q *RangeQuery) Metric
	FetchRateRange(metricName string, labels []string, grouping string, q *RangeQuery) Metric

	FetchNodeMetricValue(serviceAlias string) (int, int)
}

func (in *Client) FetchRange(metricName, labels, grouping, aggregator string, q *RangeQuery) Metric {
	query := fmt.Sprintf("%s(%s%s)", aggregator, metricName, labels)
	if grouping != "" {
		query += fmt.Sprintf(" by (%s)", grouping)
	}
	return fetchRange(in.ctx, in.api, query, q.Range)
}

func (in *Client) FetchRateRange(metricName string, labels []string, grouping string, q *RangeQuery) Metric {
	return fetchRateRange(in.ctx, in.api, metricName, labels, grouping, q)
}

func (in *Client) FetchHistogramRange(metricName, labels, grouping string, q *RangeQuery) Histogram {
	return fetchHistogramRange(in.ctx, in.api, metricName, labels, grouping, q)
}

func (in *Client) FetchHistogramValues(metricName, labels, grouping, rateInterval string, avg bool, quantiles []string, queryTime time.Time) (map[string]model.Vector, error) {
	return fetchHistogramValues(in.ctx, in.api, metricName, labels, grouping, rateInterval, avg, quantiles, queryTime)
}
func (in *Client) GenGraph(ns string, app string) (g graph.Graph, err error) {
	nodeMap := make(map[string]*graph.Node)
	nodeMap["istio-istio"] = &graph.Node{
		Workload: "istio-istio",
		NodeType: "gateway",
	}

	valueHttp, _, err := in.api.Query(in.ctx, fmt.Sprintf(`istio_requests_total{namespace="%s",rainbond_app="%s"}`, ns, app), time.Now())
	if err != nil {
		return g, err
	}
	graph.AddNode(valueHttp, nodeMap)

	valueTcp, _, err := in.api.Query(in.ctx, fmt.Sprintf(`istio_tcp_sent_bytes_total{namespace="%s",rainbond_app="%s"}`, ns, app), time.Now())
	if err != nil {
		return g, err
	}
	graph.AddNode(valueTcp, nodeMap)
	g.Nodes = nodeMap

	var edges []*graph.Edge
	edges = graph.AddEdge(valueHttp, edges, nodeMap)
	g.Edges = graph.AddEdge(valueTcp, edges, nodeMap)

	for s := range nodeMap {
		if string(s) == "istio-istio" || nodeMap[s].Protocol != "http" {
			continue
		}
		t, r := in.FetchNodeMetricValue(nodeMap[s].ServiceAlias)
		nodeMap[s].RequestRate = r
		nodeMap[s].RequestTime = t
	}
	for _, e := range edges {
		e.RequestTime = nodeMap[e.Dist].RequestTime
		e.RequestRate = nodeMap[e.Dist].RequestRate
	}
	return g, nil
}

func (in *Client) FetchNodeMetricValue(serviceAlias string) (string, string) {
	RequestTimeValue := "0"
	RequestRateValue := "0"
	rateInterval := "30s"
	lb := NewMetricsLabelsBuilder("source")
	lb.Add("service_alias", serviceAlias)

	RequestRate, _, err := in.api.Query(in.ctx, fmt.Sprintf(`sum(irate(istio_requests_total%s[%s]))`, lb.Build(), rateInterval), time.Now())
	if err == nil {
		RequestRateValue = RequestRate.(model.Vector)[0].Value.String()
	}

	query := fmt.Sprintf("sum(irate(%s_sum%s[%s])) / sum(irate(%s_count%s[%s]))",
		"istio_request_duration_milliseconds", lb.Build(), rateInterval, "istio_request_duration_milliseconds", lb.Build(), rateInterval)

	RequestTime, _, err := in.api.Query(in.ctx, query, time.Now())
	if err == nil {
		RequestTimeValue = RequestTime.(model.Vector)[0].Value.String()
	}

	//fmt.Println(lb.BuildForErrors()[0])
	//query2 := fmt.Sprintf("sum(increase(istio_requests_total%s[%s])) / sum(increase(istio_requests_total%s[%s]))", lb.BuildForErrors()[0], rateInterval, lb.Build(), rateInterval)
	//fmt.Println(query2)
	//RequestError, _, err := in.api.Query(in.ctx, query2, time.Now())
	//marshal, err := json.Marshal(RequestError)
	//if err != nil {
	//	return MetricValue{}
	//}
	//fmt.Println(string(marshal))
	//
	//if err == nil {
	//	RequestErrorRateValue = RequestError.(model.Vector)[0].Value.String()
	//} else {
	//	fmt.Println(err.Error())
	//}
	fmt.Println(formatStringToInt(RequestTimeValue))
	return formatStringToInt(RequestTimeValue), formatStringToInt(RequestRateValue)

}
