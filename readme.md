# latency-monitor

Exchange small UDP datagrams between peers to measure latency between them, and
report the statistics as prometheus metrics.

## TL;DR

```shell
latency-monitor serve \
  --transponder-listen-address 127.0.0.1:32123
```

```shell
curl -sS 127.0.0.1:8080/metrics | grep -v "^#.*$" | sort -u | grep "latency_monitor"
```

```text
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="+Inf"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="1"} 0
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="1.5"} 0
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="1000"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="115.5"} 0
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="115478"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="13.5"} 0
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="13335"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="1540"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="177828"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="178"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="1e+06"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="2.5"} 0
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="20.5"} 0
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="20535.5"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="2371.5"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="273842"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="274"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="3.5"} 0
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="31.5"} 0
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="31623"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="3651.5"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="421.5"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="421696.5"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="48.5"} 0
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="48697"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="5.5"} 0
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="5623.5"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="649.5"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="649381.5"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="74989.5"} 1
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="75"} 0
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="8.5"} 0
latency_monitor_forward_trip_latency_microseconds_bucket{peer="localhost",le="8659.5"} 1
latency_monitor_forward_trip_latency_microseconds_count{peer="localhost"} 1
latency_monitor_forward_trip_latency_microseconds_sum{peer="localhost"} 162
latency_monitor_probe_returned_count_total{peer="localhost"} 1
latency_monitor_probe_sent_count_total{peer="localhost"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="+Inf"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="1"} 0
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="1.5"} 0
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="1000"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="115.5"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="115478"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="13.5"} 0
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="13335"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="1540"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="177828"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="178"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="1e+06"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="2.5"} 0
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="20.5"} 0
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="20535.5"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="2371.5"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="273842"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="274"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="3.5"} 0
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="31.5"} 0
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="31623"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="3651.5"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="421.5"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="421696.5"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="48.5"} 0
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="48697"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="5.5"} 0
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="5623.5"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="649.5"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="649381.5"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="74989.5"} 1
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="75"} 0
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="8.5"} 0
latency_monitor_return_trip_latency_microseconds_bucket{peer="localhost",le="8659.5"} 1
latency_monitor_return_trip_latency_microseconds_count{peer="localhost"} 1
latency_monitor_return_trip_latency_microseconds_sum{peer="localhost"} 84
```

>
> Note: There is a specially-treated peer name `localhost` that can be used to
>       send probes to self.  This can be used to adjust the remote peers's
>       latencies for locally incurred implicit ones (that is, by subtracting
>       local average latency from all remote ones).
>
