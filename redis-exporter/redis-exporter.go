package main

import (
	"encoding/json"
	"flag"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type Exporter struct {
	newStatusMetric map[string]*prometheus.Desc
}

type Cpu []struct {
	Timestamp  int64   `json:"timestamp"`
	UserID     string  `json:"userId"`
	InstanceID string  `json:"instanceId"`
	Maximum    float64 `json:"Maximum"`
	Minimum    float64 `json:"Minimum"`
	Average    float64 `json:"Average"`
}

func newECSMetric(metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName("aliyun", "redis", metricName),
		docString, labels, nil,
	)
}

func newExporter() *Exporter {
	return &Exporter{
		newStatusMetric: map[string]*prometheus.Desc{
			"ConnectionUsage":  newECSMetric("ConnectionUsage", "ConnectionUsage", []string{"id"}),
			"CpuUsage":         newECSMetric("CpuUsage", "CpuUsage", []string{"id"}),
			"FailedCount":      newECSMetric("FailedCount", "FailedCount", []string{"id"}),
			"IntranetIn":       newECSMetric("IntranetIn", "IntranetIn", []string{"id"}),
			"IntranetInRatio":  newECSMetric("IntranetInRatio", "IntranetInRatio", []string{"id"}),
			"IntranetOut":      newECSMetric("IntranetOut", "IntranetOut", []string{"id"}),
			"IntranetOutRatio": newECSMetric("IntranetOutRatio", "IntranetOutRatio", []string{"id"}),
			"MemoryUsage":      newECSMetric("MemoryUsage", "MemoryUsage", []string{"id"}),
			"UsedConnection":   newECSMetric("UsedConnection", "UsedConnection", []string{"id"}),
			"UsedMemory":       newECSMetric("UsedMemory", "UsedMemory", []string{"id"}),
			"UsedQPS":          newECSMetric("UsedQPS", "UsedQPS", []string{"id"}),
		},
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.newStatusMetric {
		ch <- m
	}
}

func (e Exporter) Collect(ch chan<- prometheus.Metric) {
	for metric, desc := range e.newStatusMetric {
		request := cms.CreateDescribeMetricLastRequest()
		request.Scheme = "https"
		request.MetricName = metric
		request.Namespace = "acs_kvstore"
		request.AcceptFormat = "json"
		client, _ := cms.NewClientWithAccessKey("cn-hangzhou", "secretid", "secretkey")
		response, err := client.DescribeMetricLast(request)
		if err != nil {
			continue
		}
		var user Cpu
		json.Unmarshal([]byte(string(response.Datapoints)), &user)
		for _, value := range user {
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value.Average, value.InstanceID)
		}
	}
}

var (
	listenAddress   = flag.String("telemetry.address", ":8025", "Address on which to expose metrics.")
	metricsEndpoint = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
	//insecure        = flag.Bool("insecure", true, "Ignore server certificate if using https")
)

func main() {
	flag.Parse()
	exporter := newExporter()
	prometheus.MustRegister(exporter)
	prometheus.Unregister(prometheus.NewGoCollector())

	http.Handle(*metricsEndpoint, promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
