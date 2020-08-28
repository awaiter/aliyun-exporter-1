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
	Port       string  `json:"port"`
	Protocol   string  `json:"protocol"`
	Vip        string  `json:"vip"`
	Maximum    float64 `json:"Maximum"`
	Minimum    float64 `json:"Minimum"`
	Average    float64 `json:"Average"`
}

func newECSMetric(metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName("aliyun", "slb", metricName),
		docString, labels, nil,
	)
}

func newExporter() *Exporter {
	return &Exporter{
		newStatusMetric: map[string]*prometheus.Desc{
			"ActiveConnection":                 newECSMetric("ActiveConnection", "ActiveConnection", []string{"id", "protocol", "port", "vip"}),
			"DropConnection":                   newECSMetric("DropConnection", "DropConnection", []string{"id", "protocol", "port", "vip"}),
			"DropPackerRX":                     newECSMetric("DropPackerRX", "DropPackerRX", []string{"id", "protocol", "port", "vip"}),
			"DropPackerTX":                     newECSMetric("DropPackerTX", "DropPackerTX", []string{"id", "protocol", "port", "vip"}),
			"DropTrafficRX":                    newECSMetric("DropTrafficRX", "DropTrafficRX", []string{"id", "protocol", "port", "vip"}),
			"DropTrafficTX":                    newECSMetric("DropTrafficTX", "DropTrafficTX", []string{"id", "protocol", "port", "vip"}),
			"HeathyServerCount":                newECSMetric("HeathyServerCount", "HeathyServerCount", []string{"id", "protocol", "port", "vip"}),
			"InactiveConnection":               newECSMetric("InactiveConnection", "InactiveConnection", []string{"id", "protocol", "port", "vip"}),
			"InstanceActiveConnection":         newECSMetric("InstanceActiveConnection", "InstanceActiveConnection", []string{"id", "protocol", "port", "vip"}),
			"InstanceDropConnection":           newECSMetric("InstanceDropConnection", "InstanceDropConnection", []string{"id", "protocol", "port", "vip"}),
			"InstanceDropPacketRX":             newECSMetric("InstanceDropPacketRX", "InstanceDropPacketRX", []string{"id", "protocol", "port", "vip"}),
			"InstanceDropPacketTX":             newECSMetric("InstanceDropPacketTX", "InstanceDropPacketTX", []string{"id", "protocol", "port", "vip"}),
			"InstanceDropTrafficRX":            newECSMetric("InstanceDropTrafficRX", "InstanceDropTrafficRX", []string{"id", "protocol", "port", "vip"}),
			"InstanceDropTrafficTX":            newECSMetric("InstanceDropTrafficTX", "InstanceDropTrafficTX", []string{"id", "protocol", "port", "vip"}),
			"InstanceInactiveConnection":       newECSMetric("InstanceInactiveConnection", "InstanceInactiveConnection", []string{"id", "protocol", "port", "vip"}),
			"InstanceMaxConnection":            newECSMetric("InstanceMaxConnection", "InstanceMaxConnection", []string{"id", "protocol", "port", "vip"}),
			"InstanceMaxConnectionUtilization": newECSMetric("InstanceMaxConnectionUtilization", "InstanceMaxConnectionUtilization", []string{"id", "protocol", "port", "vip"}),
			"InstanceNewConnection":            newECSMetric("InstanceNewConnection", "InstanceNewConnection", []string{"id", "protocol", "port", "vip"}),
			"InstanceNewConnectionUtilization": newECSMetric("InstanceNewConnectionUtilization", "InstanceNewConnectionUtilization", []string{"id", "protocol", "port", "vip"}),
			"InstancePacketRX":                 newECSMetric("InstancePacketRX", "InstancePacketRX", []string{"id", "protocol", "port", "vip"}),
			"InstancePacketTX":                 newECSMetric("InstancePacketTX", "InstancePacketTX", []string{"id", "protocol", "port", "vip"}),
			"InstanceQps":                      newECSMetric("InstanceQps", "InstanceQps", []string{"id", "protocol", "port", "vip"}),
			"InstanceQpsUtilization":           newECSMetric("InstanceQpsUtilization", "InstanceQpsUtilization", []string{"id", "protocol", "port", "vip"}),
			"InstanceRt":                       newECSMetric("InstanceRt", "InstanceRt", []string{"id", "protocol", "port", "vip"}),
			"InstanceStatusCode2xx":            newECSMetric("InstanceStatusCode2xx", "InstanceStatusCode2xx", []string{"id", "protocol", "port", "vip"}),
			"InstanceStatusCode3xx":            newECSMetric("InstanceStatusCode3xx", "InstanceStatusCode3xx", []string{"id", "protocol", "port", "vip"}),
			"InstanceStatusCode4xx":            newECSMetric("InstanceStatusCode4xx", "InstanceStatusCode4xx", []string{"id", "protocol", "port", "vip"}),
			"InstanceStatusCode5xx":            newECSMetric("InstanceStatusCode5xx", "InstanceStatusCode5xx", []string{"id", "protocol", "port", "vip"}),
			"InstanceStatusCodeOther":          newECSMetric("InstanceStatusCodeOther", "InstanceStatusCodeOther", []string{"id", "protocol", "port", "vip"}),
			"InstanceTrafficRX":                newECSMetric("InstanceTrafficRX", "InstanceTrafficRX", []string{"id", "protocol", "port", "vip"}),
			"InstanceTrafficTX":                newECSMetric("InstanceTrafficTX", "InstanceTrafficTX", []string{"id", "protocol", "port", "vip"}),
			"InstanceUpstreamCode4xx":          newECSMetric("InstanceUpstreamCode4xx", "InstanceUpstreamCode4xx", []string{"id", "protocol", "port", "vip"}),
			"InstanceUpstreamCode5xx":          newECSMetric("InstanceUpstreamCode5xx", "InstanceUpstreamCode5xx", []string{"id", "protocol", "port", "vip"}),
			"InstanceUpstreamRt":               newECSMetric("InstanceUpstreamRt", "InstanceUpstreamRt", []string{"id", "protocol", "port", "vip"}),
			"MaxConnection":                    newECSMetric("MaxConnection", "MaxConnection", []string{"id", "protocol", "port", "vip"}),
			"NewConnection":                    newECSMetric("NewConnection", "NewConnection", []string{"id", "protocol", "port", "vip"}),
			"PacketRX":                         newECSMetric("PacketRX", "PacketRX", []string{"id", "protocol", "port", "vip"}),
			"PacketTX":                         newECSMetric("PacketTX", "PacketTX", []string{"id", "protocol", "port", "vip"}),
			"Qps":                              newECSMetric("Qps", "Qps", []string{"id", "protocol", "port", "vip"}),
			"Rt":                               newECSMetric("Rt", "Rt", []string{"id", "protocol", "port", "vip"}),
			"StatusCode2xx":                    newECSMetric("StatusCode2xx", "StatusCode2xx", []string{"id", "protocol", "port", "vip"}),
			"StatusCode3xx":                    newECSMetric("StatusCode3xx", "StatusCode3xx", []string{"id", "protocol", "port", "vip"}),
			"StatusCode4xx":                    newECSMetric("StatusCode4xx", "StatusCode4xx", []string{"id", "protocol", "port", "vip"}),
			"StatusCode5xx":                    newECSMetric("StatusCode5xx", "StatusCode5xx", []string{"id", "protocol", "port", "vip"}),
			"StatusCodeOther":                  newECSMetric("StatusCodeOther", "StatusCodeOther", []string{"id", "protocol", "port", "vip"}),
			"TrafficRXNew":                     newECSMetric("TrafficRXNew", "TrafficRXNew", []string{"id", "protocol", "port", "vip"}),
			"TrafficTXNew":                     newECSMetric("TrafficTXNew", "TrafficTXNew", []string{"id", "protocol", "port", "vip"}),
			"UnhealthyServerCount":             newECSMetric("UnhealthyServerCount", "UnhealthyServerCount", []string{"id", "protocol", "port", "vip"}),
			"UpstreamCode4xx":                  newECSMetric("UpstreamCode4xx", "UpstreamCode4xx", []string{"id", "protocol", "port", "vip"}),
			"UpstreamCode5xx":                  newECSMetric("UpstreamCode5xx", "UpstreamCode5xx", []string{"id", "protocol", "port", "vip"}),
			"UpstreamRt":                       newECSMetric("UpstreamRt", "UpstreamRt", []string{"id", "protocol", "port", "vip"}),
			"GroupTotalTrafficRX":              newECSMetric("GroupTotalTrafficRX", "GroupTotalTrafficRX", []string{"id", "protocol", "port", "vip"}),
			"GroupTotalTrafficTX":              newECSMetric("GroupTotalTrafficTX", "GroupTotalTrafficTX", []string{"id", "protocol", "port", "vip"}),
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
		request.Namespace = "acs_slb_dashboard"
		request.AcceptFormat = "json"
		client, _ := cms.NewClientWithAccessKey("cn-hangzhou", "secretid", "secretkey")
		response, err := client.DescribeMetricLast(request)
		if err != nil {
			continue
		}
		var user Cpu
		json.Unmarshal([]byte(string(response.Datapoints)), &user)
		for _, value := range user {
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value.Average, value.InstanceID, value.Protocol, value.Port, value.Vip)
		}
	}
}

var (
	listenAddress   = flag.String("telemetry.address", ":8026", "Address on which to expose metrics.")
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
