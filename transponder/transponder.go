package transponder

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/flashbots/latency-monitor/config"
	"github.com/flashbots/latency-monitor/logutils"
	"github.com/flashbots/latency-monitor/types"
	"go.uber.org/zap"
)

type Transponder struct {
	Receive Receive

	ip   net.IP
	port int

	conn         *net.UDPConn
	mx           sync.Mutex
	shuttingDown bool
}

type Receive = func(t *Transponder, b []byte, addr *net.UDPAddr)

var (
	ErrAlreadyServing         = errors.New("probe-responder is already serving")
	ErrMalformedListenAddress = errors.New("malformed listen address")
)

func New(cfg *config.Transponder) (*Transponder, error) {
	parts := strings.Split(cfg.ListenAddress, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: %s",
			ErrMalformedListenAddress, cfg.ListenAddress,
		)
	}

	ip := net.ParseIP(parts[0])
	if ip == nil {
		return nil, fmt.Errorf("%w: %s",
			ErrMalformedListenAddress, cfg.ListenAddress,
		)
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("%w: %w",
			ErrMalformedListenAddress, err,
		)
	}

	return &Transponder{
		ip:   ip,
		port: port,
	}, nil
}

func (t *Transponder) Shutdown(ctx context.Context) error {
	t.mx.Lock()
	defer t.mx.Unlock()

	if t.conn != nil {
		t.shuttingDown = true
		return t.conn.Close()
	}
	return nil
}

func (t *Transponder) Run(ctx context.Context) error {
	if err := t.setupConnection(); err != nil {
		return err
	}

	l := logutils.LoggerFromContext(ctx)

	buf := make([]byte, types.ProbeSize()) // must be larger than encoded Probe size

	for {
		length, addr, err := t.conn.ReadFromUDP(buf)
		if err != nil {
			if t.shuttingDown {
				return nil
			}
			l.Error("Error while reading UDP",
				zap.Error(err),
			)
			return err
		}

		t.Receive(t, buf[:length], addr)
	}
}

func (t *Transponder) IsRunning() bool {
	t.mx.Lock()
	defer t.mx.Unlock()

	return !t.shuttingDown && t.conn != nil
}

func (t *Transponder) setupConnection() error {
	t.mx.Lock()
	defer t.mx.Unlock()

	if t.conn != nil {
		return ErrAlreadyServing
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   t.ip,
		Port: t.port,
	})
	if err != nil {
		return err
	}

	t.conn = conn
	return nil
}

func (t *Transponder) Send(data []byte, addr *net.UDPAddr, onError func(error)) {
	if _, err := t.conn.WriteToUDP(data, addr); err != nil {
		onError(err)
	}
}
