package controller

import (
	"github.com/gin-gonic/gin"
	"istio-server/prometheus"
	"time"
)

func fillQueryRange(p *QueryMetricParam) *prometheus.RangeQuery {
	var q = &prometheus.RangeQuery{}
	q.FillDefaults()
	if p.Step != 0 {
		q.Step = time.Second * time.Duration(p.Step)
	}

	if p.EndTime == 0 {
		p.EndTime = time.Now().Unix()
	}
	if p.StartTime == 0 {
		p.StartTime = p.EndTime - 900
	}
	//q.RateInterval = strconv.FormatInt(p.EndTime-p.StartTime, 10) + "s"
	q.Start = time.Unix(p.StartTime, 0)
	q.End = time.Unix(p.EndTime, 0)

	return q
}

func fillQueryLabel(c *gin.Context) *prometheus.MetricsLabelsBuilder {
	lb := prometheus.NewMetricsLabelsBuilder("source")
	lb.Add("service_alias", c.Param("name"))
	//lb.Add("rainbond_app", c.Param("ns"))
	return lb
}
