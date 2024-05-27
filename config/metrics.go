package config

type Metrics struct {
	ListenAddress string `yaml:"metrics_listen_address"`

	LatencyBucketsCount int `yaml:"metrics_latency_buckets_count"`
	MaxLatencyUs        int `yaml:"metrics_max_latency_us"`

	MetricsVersion string `yaml:"metrics_version"`
}
