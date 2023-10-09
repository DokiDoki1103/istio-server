package prometheus

import (
	"fmt"
	"strconv"
	"time"

	prom_v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type RangeQuery struct {
	prom_v1.Range
	RateInterval string
	RateFunc     string
	Quantiles    []string
	Avg          bool
	ByLabels     []string
}

func (q *RangeQuery) FillDefaults() {
	q.End = time.Now()
	q.Start = q.End.Add(-30 * time.Minute)
	q.Step = 15 * time.Second
	q.RateInterval = "1m"
	q.RateFunc = "rate"
	q.Avg = true
}

type Metrics struct {
	Metrics    map[string]*Metric   `json:"metrics"`
	Histograms map[string]Histogram `json:"histograms"`
}

type Metric struct {
	Matrix model.Matrix `json:"matrix"`
	Err    error        `json:"-"`
}

type Histogram = map[string]Metric

func formatStringToInt(floatStr string) string {

	if floatStr == "NaN" {
		return "0"
	}
	// 将浮点数字符串转换为浮点数
	floatValue, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		fmt.Println("无法将字符串转换为浮点数:", err)
		return "0"
	}

	formattedStr := strconv.FormatFloat(floatValue, 'f', 2, 64)
	return formattedStr
}
