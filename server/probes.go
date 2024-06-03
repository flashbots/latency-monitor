package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"

	"time"

	"github.com/flashbots/latency-monitor/logutils"
	"github.com/flashbots/latency-monitor/metrics"
	"github.com/flashbots/latency-monitor/transponder"
	"github.com/flashbots/latency-monitor/types"
	otelattr "go.opentelemetry.io/otel/attribute"
	otelapi "go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

var (
	ErrUnexpectedDstUUIDOnReturn = errors.New("unexpected destination uuid on probe's return")
	ErrUnexpectedSrcDstUUIDs     = errors.New("source uuid is not us, but non-zero destination uuid")
)

func (s *Server) sendProbes(ctx context.Context, t *transponder.Transponder) {
	l := logutils.LoggerFromContext(ctx)

	for peerUUID, peer := range s.peers {
		addr, err := peer.UDPAddress()
		if err != nil {
			metrics.CounterFailedProbeSend.Add(ctx, 1, s.labels, otelapi.WithAttributes(
				otelattr.String("error_type", reflect.TypeOf(err).String()),
			))
			l.Error("Failed to send a probe",
				zap.Error(err),
				zap.String("peer", peer.Name()),
			)
			continue
		}

		p := types.Probe{
			Sequence: peer.Sequence(),
			SrcUUID:  s.uuid,
			DstUUID:  peerUUID,
		}
		p.SrcTimestamp = time.Now()

		b, err := p.MarshalBinary()
		if err != nil {
			metrics.CounterFailedProbeSend.Add(ctx, 1, s.labels, otelapi.WithAttributes(
				otelattr.String("error_type", reflect.TypeOf(err).String()),
			))
			l.Error("Failed to prepare a probe",
				zap.Error(err),
			)
			continue
		}

		t.Send(b, addr, func(err error) {
			metrics.CounterFailedProbeSend.Add(ctx, 1, s.labels, otelapi.WithAttributes(
				otelattr.String("error_type", reflect.TypeOf(err).String()),
			))
			l.Error("Failed to send a probe",
				zap.Error(err),
			)
		})

		metrics.CountProbeSent.Add(ctx, 1, s.labels, otelapi.WithAttributes(
			otelattr.String("peer", peer.Name()),
		))
		l.Debug("Sent a probe",
			zap.String("name", peer.Name()),
		)
	}
}

func (s *Server) receiveProbes(ctx context.Context) transponder.Receive {
	l := logutils.LoggerFromContext(ctx)

	return func(t *transponder.Transponder, input []byte, source *net.UDPAddr) {
		ts := time.Now()

		p := types.Probe{}
		if err := p.UnmarshalBinary(input); err != nil {
			metrics.CounterInvalidProbeReceived.Add(ctx, 1, s.labels, otelapi.WithAttributes(
				otelattr.String("error_type", reflect.TypeOf(err).String()),
			))
			l.Error("Invalid probe",
				zap.Error(err),
				zap.String("source", source.String()),
				zap.ByteString("payload", input),
			)
			return
		}

		switch {
		case p.DstTimestamp.IsZero(): // reply to the others' probes
			p.DstTimestamp = ts
			output, err := p.MarshalBinary()
			if err != nil {
				metrics.CounterFailedProbeRespond.Add(ctx, 1, s.labels, otelapi.WithAttributes(
					otelattr.String("error_type", reflect.TypeOf(err).String()),
				))
				l.Error("Failed to prepare response to a probe",
					zap.Error(err),
				)
				return
			}

			go func() {
				t.Send(output, source, func(err error) {
					metrics.CounterFailedProbeRespond.Add(ctx, 1, s.labels, otelapi.WithAttributes(
						otelattr.String("error_type", reflect.TypeOf(err).String()),
					))
					l.Error("Failed to respond to a probe",
						zap.Error(err),
					)
				})
			}()

		case p.SrcUUID == s.uuid: // handle our own (returned) probes
			peer, known := s.peers[p.DstUUID]
			if !known {
				err := fmt.Errorf("%w: %s",
					ErrUnexpectedDstUUIDOnReturn, p.DstUUID.String(),
				)
				metrics.CounterInvalidProbeReceived.Add(ctx, 1, s.labels, otelapi.WithAttributes(
					otelattr.String("error_type", reflect.TypeOf(err).String()),
				))
				l.Error("Invalid return probe",
					zap.Error(err),
					zap.String("source", source.String()),
				)
				return
			}

			forwardLatency := float64(p.DstTimestamp.Sub(p.SrcTimestamp).Microseconds())
			metrics.HistogramLatencyForwardTrip.Record(ctx, forwardLatency, s.labels, otelapi.WithAttributes(
				otelattr.String("peer", peer.Name()),
			))

			returnLatency := float64(ts.Sub(p.DstTimestamp).Microseconds())
			metrics.HistogramLatencyReturnTrip.Record(ctx, returnLatency, s.labels, otelapi.WithAttributes(
				otelattr.String("peer", peer.Name()),
			))

			metrics.CountProbeReturned.Add(ctx, 1, s.labels, otelapi.WithAttributes(
				otelattr.String("peer", peer.Name()),
			))
			l.Debug("Received a return probe",
				zap.Float64("forward_latency_ms", forwardLatency),
				zap.Float64("return_latency_ms", returnLatency),
				zap.String("name", peer.Name()),
			)

			return

		default: // handle mismatching probes
			err := fmt.Errorf("%w: source %s, destination %s",
				ErrUnexpectedSrcDstUUIDs, p.SrcUUID.String(), p.DstUUID.String(),
			)
			metrics.CounterInvalidProbeReceived.Add(ctx, 1, s.labels, otelapi.WithAttributes(
				otelattr.String("error_type", reflect.TypeOf(err).String()),
			))
			l.Error("Invalid probe",
				zap.Error(err),
			)
			return
		}
	}
}
