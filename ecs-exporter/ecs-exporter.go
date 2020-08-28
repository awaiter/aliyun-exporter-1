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
	UserID     string  `json:"userId,omitempty"`
	InstanceID string  `json:"instanceId"`
	Device     string  `json:"device"`
	Hostname   string  `json:"hostname"`
	IP         string  `json:"IP"`
	Sum        float64 `json:"Sum"`
	Maximum    float64 `json:"Maximum"`
	Average    float64 `json:"Average"`
	Minimum    float64 `json:"Minimum"`
	State      string  `json:"state"`
	Diskname   string  `json:"diskname"`
}

func newECSMetric(metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName("aliyun", "ecs", metricName),
		docString, labels, nil,
	)
}

func newExporter() *Exporter {
	return &Exporter{
		newStatusMetric: map[string]*prometheus.Desc{
			"cpu_total":                 newECSMetric("cpu_total", "cpu_total", []string{"id"}),
			"cpu_idle":                  newECSMetric("cpu_idle", "cpu_idle", []string{"id"}),
			"cpu_other":                 newECSMetric("cpu_other", "cpu_other", []string{"id"}),
			"cpu_system":                newECSMetric("cpu_system", "cpu_system", []string{"id"}),
			"cpu_user":                  newECSMetric("cpu_user", "cpu_user", []string{"id"}),
			"cpu_wait":                  newECSMetric("cpu_wait", "cpu_wait", []string{"id"}),
			"disk_readbytes":            newECSMetric("disk_readbytes", "disk_readbytes", []string{"id", "device"}),
			"disk_readiops":             newECSMetric("disk_readiops", "disk_readiops", []string{"id", "device"}),
			"disk_writebytes":           newECSMetric("disk_writebytes", "disk_writebytes", []string{"id", "device"}),
			"disk_writeiops":            newECSMetric("disk_writeiops", "disk_writeiops", []string{"id", "device"}),
			"diskusage_free":            newECSMetric("diskusage_free", "diskusage_free", []string{"id", "device", "diskname", "hostname"}),
			"diskusage_avail":           newECSMetric("diskusage_avail", "diskusage_avail", []string{"id", "device", "diskname", "hostname"}),
			"diskusage_total":           newECSMetric("diskusage_total", "diskusage_total", []string{"id", "device", "diskname", "hostname"}),
			"diskusage_used":            newECSMetric("diskusage_used", "diskusage_used", []string{"id", "device", "diskname", "hostname"}),
			"diskusage_utilization":     newECSMetric("diskusage_utilization", "diskusage_utilization", []string{"id", "device", "diskname", "hostname"}),
			"fs_inodeutilization":       newECSMetric("fs_inodeutilization", "fs_inodeutilization", []string{"id", "device", "diskname", "hostname"}),
			"load_15m":                  newECSMetric("load_15m", "load_15m", []string{"id"}),
			"load_1m":                   newECSMetric("load_1m", "load_1m", []string{"id"}),
			"load_5m":                   newECSMetric("load_5m", "load_5m", []string{"id"}),
			"memory_freespace":          newECSMetric("memory_freespace", "memory_freespace", []string{"id"}),
			"memory_freeutilization":    newECSMetric("memory_freeutilization", "memory_freeutilization", []string{"id"}),
			"memory_totalspace":         newECSMetric("memory_totalspace", "memory_totalspace", []string{"id"}),
			"memory_usedspace":          newECSMetric("memory_usedspace", "memory_usedspace", []string{"id"}),
			"memory_usedutilization":    newECSMetric("memory_usedutilization", "memory_usedutilization", []string{"id"}),
			"net_tcpconnection":         newECSMetric("net_tcpconnection", "net_tcpconnection", []string{"id", "state"}),
			"networkin_errorpackages":   newECSMetric("networkin_errorpackages", "networkin_errorpackages", []string{"id", "Device"}),
			"networkin_packages":        newECSMetric("networkin_packages", "networkin_packages", []string{"id", "device", "interface"}),
			"networkin_packages_total":  newECSMetric("networkin_packages_total", "networkin_packages_total", []string{"id", "device", "interface"}),
			"networkin_rate":            newECSMetric("networkin_rate", "networkin_rate", []string{"id", "device", "interface"}),
			"networkout_errorpackages":  newECSMetric("networkout_errorpackages", "networkout_errorpackages", []string{"id", "Device"}),
			"networkout_packages":       newECSMetric("networkout_packages", "networkout_packages", []string{"id", "device", "interface"}),
			"networkout_packages_total": newECSMetric("networkout_packages_total", "networkout_packages_total", []string{"id", "device", "interface"}),
			"networkout_rate":           newECSMetric("networkout_rate", "networkout_rate", []string{"id", "device", "interface"}),
			"process_number":            newECSMetric("process_number", "process_number", []string{"id"}),
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
		request.Namespace = "acs_ecs_dashboard"
		request.AcceptFormat = "json"
		client, _ := cms.NewClientWithAccessKey("cn-hangzhou", "secretid", "secretkey")
		response, err := client.DescribeMetricLast(request)
		if err != nil {
			continue
		}
		var user Cpu
		json.Unmarshal([]byte(string(response.Datapoints)), &user)
		for _, value := range user {
			if metric == "cpu_total" || metric == "cpu_idle" || metric == "cpu_other" || metric == "cpu_system" || metric == "cpu_user" || metric == "cpu_wait" || metric == "load_15m" || metric == "load_1m" || metric == "load_5m" || metric == "memory_freespace" || metric == "memory_freeutilization" || metric == "memory_totalspace" || metric == "memory_usedspace" || metric == "memory_usedutilization" {
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value.Average, value.InstanceID)
			} else if metric == "net_tcpconnection" {
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value.Average, value.InstanceID, value.State)
			} else if metric == "networkin_errorpackages" || metric == "networkout_errorpackages" {
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value.Average, value.InstanceID, value.Device)
			} else if metric == "networkin_packages" || metric == "networkout_packages" || metric == "networkin_rate" || metric == "networkout_rate" {
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value.Average, value.InstanceID, value.Device, value.IP)
				if metric == "networkin_packages" {
					ch <- prometheus.MustNewConstMetric(e.newStatusMetric["networkin_packages_total"], prometheus.GaugeValue, value.Sum, value.InstanceID, value.Device, value.IP)
				} else if metric == "networkout_packages" {
					ch <- prometheus.MustNewConstMetric(e.newStatusMetric["networkout_packages_total"], prometheus.GaugeValue, value.Sum, value.InstanceID, value.Device, value.IP)
				}
			} else if metric == "fs_inodeutilization" {
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value.Average, value.InstanceID, value.Device, value.Diskname, value.Hostname)
			} else if metric == "diskusage_free" || metric == "diskusage_avail" || metric == "diskusage_total" || metric == "diskusage_used" || metric == "diskusage_utilization" {
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value.Average, value.InstanceID, value.Device, value.Diskname, value.Hostname)
			} else if metric == "disk_readbytes" || metric == "disk_writebytes" || metric == "disk_readiops" || metric == "disk_writeiops" {
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value.Average, value.InstanceID, value.Device)
			}
		}
	}
}

var (
	listenAddress   = flag.String("telemetry.address", ":8023", "Address on which to expose metrics.")
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
