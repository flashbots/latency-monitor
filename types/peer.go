package types

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Peer struct {
	name string

	host string
	port int

	udpAddress *net.UDPAddr

	sequence uint64
}

var (
	ErrPeerFailedToDecodeStringRepresentation = errors.New("failed to decode peer from its string representation")
	ErrPeerFailedToResolveIP4                 = errors.New("failed to resolve peer ip4 address")
)

func NewPeer(s string) (Peer, error) {
	p1 := strings.Split(s, "=")
	if len(p1) != 2 {
		return Peer{}, fmt.Errorf("%w: expected '=' delimiter: %s",
			ErrPeerFailedToDecodeStringRepresentation, s,
		)
	}

	p2 := strings.Split(p1[1], ":")
	if len(p2) != 2 {
		return Peer{}, fmt.Errorf("%w: expected ':' delimiter: %s",
			ErrPeerFailedToDecodeStringRepresentation, s,
		)
	}

	ip := net.ParseIP(p2[0])
	host := ""
	if ip == nil {
		if _, err := net.LookupIP(p2[0]); err != nil {
			return Peer{}, fmt.Errorf("%w: %w: %s",
				ErrPeerFailedToDecodeStringRepresentation, err, s,
			)
		}
		host = p2[0]
	}

	port, err := strconv.Atoi(p2[1])
	if err != nil {
		return Peer{}, fmt.Errorf("%w: %w",
			ErrPeerFailedToDecodeStringRepresentation, err,
		)
	}

	var udpAddress *net.UDPAddr = nil
	if ip != nil {
		udpAddress = &net.UDPAddr{
			IP:   ip,
			Port: port,
		}
	}

	return Peer{
		name: p1[0],

		host: host,
		port: port,

		udpAddress: udpAddress,
	}, nil
}

func (p Peer) Name() string {
	return p.name
}

func (p *Peer) Sequence() uint64 {
	res := p.sequence
	p.sequence += 1
	return res
}

func (p Peer) UDPAddress() (*net.UDPAddr, error) {
	if p.udpAddress != nil {
		return p.udpAddress, nil
	}

	addresses, err := net.LookupIP(p.host)
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %s",
			ErrPeerFailedToResolveIP4, err, p.host,
		)
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("%w: no ip4 found: %s",
			ErrPeerFailedToResolveIP4, p.host,
		)
	}

	for _, addr := range addresses {
		if len(addr) == 4 { // ip4
			return &net.UDPAddr{
				IP:   addr,
				Port: p.port,
			}, nil
		}
	}

	return nil, fmt.Errorf("%w: no ip4 found: %s",
		ErrPeerFailedToResolveIP4, p.host,
	)
}
