package config

import (
	"time"

	"github.com/flashbots/latency-monitor/types"
)

type Transponder struct {
	Interval      time.Duration `yaml:"transponder_interval"`
	ListenAddress string        `yaml:"transponder_listen_address"`
	Peers         []types.Peer  `yaml:"transponder_peers"`
}
