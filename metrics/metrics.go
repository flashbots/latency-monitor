package metrics

import (
	"context"
	"math"

	"github.com/flashbots/latency-monitor/config"
	"go.opentelemetry.io/otel/exporters/prometheus"
	otelapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	metricsNamespace = "latency-monitor"
)

var (
	meter               otelapi.Meter
	latencyBoundariesUs otelapi.HistogramOption

	CountProbeReturned otelapi.Int64Counter
	CountProbeSent     otelapi.Int64Counter

	CounterFailedProbeRespond   otelapi.Int64Counter
	CounterFailedProbeSend      otelapi.Int64Counter
	CounterInvalidProbeReceived otelapi.Int64Counter

	HistogramLatencyForwardTrip otelapi.Float64Histogram
	HistogramLatencyReturnTrip  otelapi.Float64Histogram
)

func Setup(ctx context.Context, cfg *config.Metrics) error {
	for _, setup := range []func(context.Context, *config.Metrics) error{
		setupMeter,               // must come first
		setupLatencyBoundariesUs, // must come second

		setupCounterProbeReturned,
		setupCounterProbeSent,

		setupCounterFailedProbeRespond,
		setupCounterInvalidProbes,
		setupCounterFailedProbeSend,

		setupHistogramLatencyForwardTrip,
		setupHistogramLatencyReturnTrip,
	} {
		if err := setup(ctx, cfg); err != nil {
			return err
		}
	}

	return nil
}

func setupMeter(ctx context.Context, cfg *config.Metrics) error {
	res, err := resource.New(ctx)

	if err != nil {
		return err
	}

	exporter, err := prometheus.New(
		prometheus.WithNamespace(metricsNamespace),
		prometheus.WithoutScopeInfo(),
	)
	if err != nil {
		return err
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(res),
	)

	meter = provider.Meter(metricsNamespace)

	return nil
}

func setupLatencyBoundariesUs(ctx context.Context, cfg *config.Metrics) error {
	latencyBoundariesUs = otelapi.WithExplicitBucketBoundaries(func() []float64 {
		base := math.Exp(math.Log(float64(cfg.MaxLatencyUs)) / (float64(cfg.LatencyBucketsCount - 1)))
		res := make([]float64, 0, cfg.LatencyBucketsCount)
		for i := 0; i < cfg.LatencyBucketsCount; i++ {
			res = append(res,
				math.Round(2*math.Pow(base, float64(i)))/2,
			)
		}
		return res
	}()...)

	return nil
}

func setupCounterProbeReturned(_ context.Context, _ *config.Metrics) error {
	counter, err := meter.Int64Counter(
		"probe_returned_count",
		otelapi.WithDescription("count of successfully returned probes"),
	)
	CountProbeReturned = counter
	if err != nil {
		return err
	}
	return nil
}

func setupCounterProbeSent(_ context.Context, _ *config.Metrics) error {
	counter, err := meter.Int64Counter(
		"probe_sent_count",
		otelapi.WithDescription("count of successfully sent probes"),
	)
	CountProbeSent = counter
	if err != nil {
		return err
	}
	return nil
}

func setupCounterFailedProbeSend(_ context.Context, _ *config.Metrics) error {
	counter, err := meter.Int64Counter(
		"failed_probe_send_count",
		otelapi.WithDescription("count of failing to send a probe"),
	)
	CounterFailedProbeSend = counter
	if err != nil {
		return err
	}
	return nil
}

func setupCounterFailedProbeRespond(_ context.Context, _ *config.Metrics) error {
	counter, err := meter.Int64Counter(
		"failed_probe_respond_count",
		otelapi.WithDescription("count of failing to respond to a probe"),
	)
	CounterFailedProbeRespond = counter
	if err != nil {
		return err
	}
	return nil
}

func setupCounterInvalidProbes(_ context.Context, _ *config.Metrics) error {
	counter, err := meter.Int64Counter(
		"invalid_probes_count",
		otelapi.WithDescription("count of receiving an invalid probe"),
	)
	CounterInvalidProbeReceived = counter
	if err != nil {
		return err
	}
	return nil
}

func setupHistogramLatencyForwardTrip(_ context.Context, _ *config.Metrics) error {
	latency, err := meter.Float64Histogram(
		"forward_trip_latency",
		otelapi.WithDescription("statistics on the latency of probes' forward-trip"),
		otelapi.WithUnit("us"),
		latencyBoundariesUs,
	)
	HistogramLatencyForwardTrip = latency
	if err != nil {
		return err
	}
	return nil
}

func setupHistogramLatencyReturnTrip(_ context.Context, _ *config.Metrics) error {
	latency, err := meter.Float64Histogram(
		"return_trip_latency",
		otelapi.WithDescription("statistics on the latency of probes' return-trip"),
		otelapi.WithUnit("us"),
		latencyBoundariesUs,
	)
	HistogramLatencyReturnTrip = latency
	if err != nil {
		return err
	}
	return nil
}
