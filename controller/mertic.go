package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/model"
	"istio-server/prometheus"
	"net/http"
	"strconv"
	"time"
)

func GetStatus(c *gin.Context) {
	c.String(200, "ok")
}

type QueryMetricParam struct {
	Step         int    `form:"step"`
	RateInterval string `form:"rateInterval"`
	StartTime    int64  `form:"startTime"`
	EndTime      int64  `form:"endTime"`
}

func GetHttpCount(c *gin.Context) {
	var client = prometheus.GetPrometheusClient()

	var p QueryMetricParam
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	queryRange := fillQueryRange(&p)
	queryRange.Step = time.Duration(p.EndTime-p.StartTime) * time.Second
	label := fillQueryLabel(c)

	metric := client.FetchRange("istio_requests_total", label.Build(), "response_code", "sum", queryRange)

	res := make(map[string]int)
	for _, m := range metric.Matrix {
		response_code := string(m.Metric["response_code"])
		if len(m.Values) == 1 {
			atoi, _ := strconv.Atoi(m.Values[0].Value.String())
			res[response_code] = atoi
		} else {
			start, _ := strconv.Atoi(m.Values[0].Value.String())
			end, _ := strconv.Atoi(m.Values[len(m.Values)-1].Value.String())
			if end-start < 0 {
				res[response_code] = end
			} else {
				res[response_code] = end - start
			}
		}
	}
	c.JSON(http.StatusOK, res)
}

func GetHttpQps(c *gin.Context) {
	var client = prometheus.GetPrometheusClient()

	var p QueryMetricParam
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	queryRange := fillQueryRange(&p)
	label := fillQueryLabel(c)

	metric := client.FetchRateRange("istio_requests_total", []string{label.Build()}, "", queryRange)

	if len(metric.Matrix) > 0 {
		c.JSON(http.StatusOK, metric.Matrix[0].Values)
	} else {
		c.JSON(http.StatusOK, []model.SamplePair{})
	}
}

func GetHttpTime(c *gin.Context) {
	var client = prometheus.GetPrometheusClient()

	var p QueryMetricParam
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	queryRange := fillQueryRange(&p)
	label := fillQueryLabel(c)

	metric := client.FetchHistogramRange("istio_request_duration_milliseconds", label.Build(), "", queryRange)

	if len(metric["avg"].Matrix) > 0 {
		c.JSON(http.StatusOK, metric["avg"].Matrix[0].Values)
	} else {
		c.JSON(http.StatusOK, []model.SamplePair{})
	}

}

func GetHttpBytes(c *gin.Context) {
	var client = prometheus.GetPrometheusClient()

	var p QueryMetricParam
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	queryRange := fillQueryRange(&p)
	label := fillQueryLabel(c)

	request := client.FetchRateRange("istio_request_bytes_sum", []string{label.Build()}, "", queryRange)
	response := client.FetchRateRange("istio_response_bytes_sum", []string{label.Build()}, "", queryRange)

	if len(request.Matrix) == 0 || len(response.Matrix) == 0 {
		c.JSON(http.StatusOK, HttpBytesRes{
			Request:  []model.SamplePair{},
			Response: []model.SamplePair{},
		})
		return
	}
	c.JSON(http.StatusOK, HttpBytesRes{
		Request:  request.Matrix[0].Values,
		Response: response.Matrix[0].Values,
	})
}

type HttpBytesRes struct {
	Request  []model.SamplePair `json:"request_bytes_sum"`
	Response []model.SamplePair `json:"response_bytes_sum"`
}

func GraphNamespaces(c *gin.Context) {
	var client = prometheus.GetPrometheusClient()
	genGraph, err := client.GenGraph(c.Query("ns"), c.Query("app"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, genGraph)
}
