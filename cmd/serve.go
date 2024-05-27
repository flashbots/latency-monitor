package main

import (
	"slices"
	"time"

	"github.com/flashbots/latency-monitor/config"
	"github.com/flashbots/latency-monitor/server"
	"github.com/flashbots/latency-monitor/types"
	"github.com/urfave/cli/v2"
)

const (
	categoryMetrics     = "METRICS:"
	categoryServer      = "SERVER:"
	categoryTransponder = "TRANSPONDER:"
)

func CommandServe(cfg *config.Config) *cli.Command {
	transponderPeers := &cli.StringSlice{}

	metricsFlags := []cli.Flag{
		&cli.StringFlag{
			Category:    categoryMetrics,
			Destination: &cfg.Metrics.ListenAddress,
			EnvVars:     []string{envPrefix + "METRICS_LISTEN_ADDRESS"},
			Name:        "metrics-listen-address",
			Usage:       "`host:port` for the metrics-server to listen on",
			Value:       "0.0.0.0:8080",
		},

		&cli.IntFlag{
			Category:    categoryMetrics,
			Destination: &cfg.Metrics.LatencyBucketsCount,
			EnvVars:     []string{envPrefix + "METRICS_LATENCY_BUCKETS_COUNT"},
			Name:        "metrics-latency-buckets-count",
			Usage:       "`count` of latency histogram buckets",
			Value:       33,
		},

		&cli.IntFlag{
			Category:    categoryMetrics,
			Destination: &cfg.Metrics.MaxLatencyUs,
			EnvVars:     []string{envPrefix + "METRICS_MAX_LATENCY"},
			Name:        "metrics-max-latency",
			Usage:       "`microseconds` value for the largest histogram latency bucket",
			Value:       1000000,
		},
	}

	transponderFlags := []cli.Flag{
		&cli.DurationFlag{
			Category:    categoryTransponder,
			Destination: &cfg.Transponder.Interval,
			EnvVars:     []string{envPrefix + "TRANSPONDER_INTERVAL"},
			Name:        "transponder-interval",
			Usage:       "`interval` at which the transponder should send its probes",
			Value:       time.Minute,
		},

		&cli.StringFlag{
			Category:    categoryTransponder,
			Destination: &cfg.Transponder.ListenAddress,
			EnvVars:     []string{envPrefix + "TRANSPONDER_LISTEN_ADDRESS"},
			Name:        "transponder-listen-port",
			Usage:       "`host:port` for the transponder to listen on",
			Value:       "0.0.0.0:32123",
		},

		&cli.StringSliceFlag{
			Category:    categoryTransponder,
			Destination: transponderPeers,
			EnvVars:     []string{envPrefix + "TRANSPONDER_PEER"},
			Name:        "transponder-peer",
			Usage:       "`name=host:port` of the transponder peer to measure the latency against",
			Value:       cli.NewStringSlice("localhost=127.0.0.1:32123"),
		},
	}

	serverFlags := []cli.Flag{
		&cli.StringFlag{
			Category:    categoryServer,
			Destination: &cfg.Server.Name,
			EnvVars:     []string{envPrefix + "SERVER_NAME"},
			Name:        "server-name",
			Usage:       "service `name` to report in prometheus metrics",
			Value:       "latency-monitor",
		},
	}

	flags := slices.Concat(
		serverFlags,
		metricsFlags,
		transponderFlags,
	)

	return &cli.Command{
		Name:  "serve",
		Usage: "run the monitor server",
		Flags: flags,

		Before: func(ctx *cli.Context) error {
			p := transponderPeers.Value()
			peers := make([]types.Peer, 0, len(p))
			for _, s := range p {
				peer, err := types.NewPeer(s)
				if err != nil {
					return err
				}
				peers = append(peers, peer)
			}
			cfg.Transponder.Peers = peers
			return nil
		},

		Action: func(_ *cli.Context) error {
			s, err := server.New(cfg)
			if err != nil {
				return err
			}
			return s.Run()
		},
	}
}
