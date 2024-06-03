package config

type Metrics struct {
	ListenAddress string `yaml:"metrics_listen_address"`

	Labels map[string]string

	LatencyBucketsCount int `yaml:"metrics_latency_buckets_count"`
	MaxLatencyUs        int `yaml:"metrics_max_latency_us"`

	Version string `yaml:"metrics_version"`
}
