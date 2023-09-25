package prometheus

import (
	"encoding/json"
	"fmt"
	"istio-server/kubernetes"
	"testing"
	"time"
)

func TestClient(*testing.T) {
	client, err := NewClient()
	if err != nil {
		fmt.Println("err")
	}

	var c = RangeQuery{}
	c.FillDefaults()
	c.Start = c.End.Add(-3 * time.Minute) //c.RateFunc = "irate"

	//lb := NewMetricsLabelsBuilder("")
	//lb.Add("name", "gr4f5ad3")

	//istio_tcp_connections_opened_total{creator="Rainbond",rainbond_app="b"}
	metric := client.FetchRange("istio_tcp_sent_bytes_total", `{rainbond_app="default",reporter="destination"}`, "", "", &c)

	marshal, err := json.Marshal(metric)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))

	//metric2 := client.FetchRateRange("istio_requests_total", []string{lb.Build()}, "response_code", &c)
	//
	//marshal2, err := json.Marshal(metric2)
	//if err != nil {
	//	return
	//}
	//fmt.Println(string(marshal2))
}

func TestDep(*testing.T) {
	_, _ = NewClient()
	var client = GetPrometheusClient()

	client.FetchNodeMetricValue("gr9efd6b")

}

func TestK8s(*testing.T) {

	client, err := kubernetes.NewClientFromConfig()
	if err != nil {
		return
	}
	password, err := client.GetMySQLPassword()
	if err != nil {
		return

	}
	fmt.Println(password)
}
