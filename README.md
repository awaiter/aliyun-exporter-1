# aliyun-exporter
prometheus exporter for aliyun

Quick Start

```
set your secretid and secretkey in the code.
client, _ := cms.NewClientWithAccessKey("cn-hangzhou", "secretid", "secretkey")

if you want to change the listen port,you can change it in the code.
listenAddress   = flag.String("telemetry.address", ":8026", "Address on which to expose metrics.")

if you want to change the endpoint, you can change it in the code.
metricsEndpoint = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics.")

thanks for useing it!!!
```

