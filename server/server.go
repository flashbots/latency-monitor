package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flashbots/latency-monitor/config"
	"github.com/flashbots/latency-monitor/httplogger"
	"github.com/flashbots/latency-monitor/logutils"
	"github.com/flashbots/latency-monitor/metrics"
	"github.com/flashbots/latency-monitor/transponder"
	"github.com/flashbots/latency-monitor/types"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	otelattr "go.opentelemetry.io/otel/attribute"
	otelapi "go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

type Server struct {
	cfg *config.Config
	log *zap.Logger

	uuid  uuid.UUID
	peers map[uuid.UUID]*types.Peer

	labels   otelapi.MeasurementOption
	location types.Location
}

func New(cfg *config.Config) (*Server, error) {
	l := zap.L()

	srvUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	labels := make([]otelattr.KeyValue, 0, len(cfg.Metrics.Labels))
	for k, v := range cfg.Metrics.Labels {
		labels = append(labels, otelattr.String(k, v))
	}

	location := types.Location{}
	copy(location[:], []byte(cfg.Metrics.Location))

	peers := make(map[uuid.UUID]*types.Peer, len(cfg.Transponder.Peers))
	for _, peer := range cfg.Transponder.Peers {
		peerUUID := srvUUID

		if peer.Name() != "localhost" {
			peerUUID, err = uuid.NewRandom()
			if err != nil {
				return nil, err
			}
		}
		peers[peerUUID] = &peer
	}

	return &Server{
		cfg: cfg,
		log: l,

		uuid:  srvUUID,
		peers: peers,

		labels:   otelapi.WithAttributeSet(otelattr.NewSet(labels...)),
		location: location,
	}, nil
}

func (s *Server) Run() error {
	l := s.log
	ctx := logutils.ContextWithLogger(context.Background(), l)

	if err := metrics.Setup(ctx, &s.cfg.Metrics); err != nil {
		return err
	}

	metricsServer, err := s.newMetricsServer(ctx)
	if err != nil {
		return err
	}

	transponder, err := transponder.New(&s.cfg.Transponder)
	if err != nil {
		return err
	}
	transponder.Receive = s.receiveProbes(ctx)

	ticker := time.NewTicker(s.cfg.Transponder.Interval)

	failure := make(chan error, 1)

	go func() { // run the transponder
		l.Info("Latency monitor transponder is going up...",
			zap.String("responder_listen_address", s.cfg.Transponder.ListenAddress),
		)
		if err := transponder.Run(ctx); err != nil {
			failure <- err
		}
		l.Info("Latency monitor transponder is down")
	}()

	go func() { // run the metrics-server
		l.Info("Latency monitor metrics-server is going up...",
			zap.String("metrics_listen_address", s.cfg.Metrics.ListenAddress),
		)
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			failure <- err
		}
		l.Info("Latency monitor metrics-server is down")
	}()

	go func() { // run the ticker
		for {
			<-ticker.C
			if !transponder.IsRunning() {
				l.Warn("Transponder is not running...")
				continue
			}
			s.sendProbes(ctx, transponder)
		}
	}()

	{ // wait until termination or internal failure
		terminator := make(chan os.Signal, 1)
		signal.Notify(terminator, os.Interrupt, syscall.SIGTERM)

		select {
		case stop := <-terminator:
			l.Info("Stop signal received; shutting down...",
				zap.String("signal", stop.String()),
			)
		case err := <-failure:
			l.Error("Internal failure; shutting down...",
				zap.Error(err),
			)
		}
	}

	{ // stop the ticker
		ticker.Stop()
	}

	{ // stop the transponder
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		if err := transponder.Shutdown(ctx); err != nil {
			l.Error("Error while shutting down latency monitor transponder",
				zap.Error(err),
			)
		}
		l.Info("Latency monitor transponder is down")
	}

	{ // stop the metrics-server
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		if err := metricsServer.Shutdown(ctx); err != nil {
			l.Error("Latency monitor metrics-server shutdown failed",
				zap.Error(err),
			)
		}
	}

	return nil
}

func (s *Server) newMetricsServer(ctx context.Context) (*http.Server, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleHealthcheck)
	mux.Handle("/metrics", promhttp.Handler())
	handler := httplogger.Middleware(s.log, mux)

	srv := &http.Server{
		Addr:              s.cfg.Metrics.ListenAddress,
		ErrorLog:          logutils.NewHttpServerErrorLogger(logutils.LoggerFromContext(ctx)),
		Handler:           handler,
		MaxHeaderBytes:    1024,
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	return srv, nil
}
