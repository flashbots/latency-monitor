package types

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Probe struct {
	Sequence     uint64
	SrcUUID      uuid.UUID
	SrcTimestamp time.Time
	DstUUID      uuid.UUID
	DstTimestamp time.Time
}

var (
	ErrProbeFailedToEncodeBinaryRepresentation = errors.New("failed to encode probe into its binary representation")
	ErrProbeFailedToDecodeBinaryRepresentation = errors.New("failed to decode probe from its binary representation")
)

func ProbeSize() int {
	return 72
}

func (p Probe) MarshalBinary() ([]byte, error) {
	rawSrcTimestamp, err := p.SrcTimestamp.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("%w: SrcTimestamp: %w",
			ErrProbeFailedToEncodeBinaryRepresentation, err,
		)
	}
	rawDstTimestamp, err := p.DstTimestamp.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("%w: DstTimestamp: %w",
			ErrProbeFailedToEncodeBinaryRepresentation, err,
		)
	}

	data := make([]byte, 72)

	binary.LittleEndian.PutUint64(data[0:8], p.Sequence) // 00..07 : 8 bytes
	copy(data[8:24], p.SrcUUID[:])                       // 08..23 : 16 bytes
	copy(data[24:39], rawSrcTimestamp)                   // 24..38 : 15 bytes
	copy(data[39:55], p.DstUUID[:])                      // 39..54 : 16 bytes
	copy(data[55:70], rawDstTimestamp)                   // 55..71 : 15 bytes

	return data, nil
}

func (p *Probe) UnmarshalBinary(data []byte) error {
	if len(data) != ProbeSize() {
		return fmt.Errorf("%w: invalid binary length: expected %d, got %d",
			ErrProbeFailedToDecodeBinaryRepresentation, ProbeSize(), len(data),
		)
	}

	srcUUID, err := uuid.FromBytes(data[8:24])
	if err != nil {
		return fmt.Errorf("%w: SrcUUID: %w",
			ErrProbeFailedToDecodeBinaryRepresentation, err,
		)
	}

	srcTimestamp := &time.Time{}
	if err := srcTimestamp.UnmarshalBinary(data[24:39]); err != nil {
		return fmt.Errorf("%w: SrcTimestamp: %w",
			ErrProbeFailedToDecodeBinaryRepresentation, err,
		)
	}

	dstUUID, err := uuid.FromBytes(data[39:55])
	if err != nil {
		return fmt.Errorf("%w: DstUUID: %w",
			ErrProbeFailedToDecodeBinaryRepresentation, err,
		)
	}

	dstTimestamp := &time.Time{}
	if err := dstTimestamp.UnmarshalBinary(data[55:70]); err != nil {
		return fmt.Errorf("%w: DstTimestamp: %w",
			ErrProbeFailedToDecodeBinaryRepresentation, err,
		)
	}

	*p = Probe{
		Sequence:     binary.LittleEndian.Uint64(data[:8]),
		SrcUUID:      srcUUID,
		SrcTimestamp: *srcTimestamp,
		DstUUID:      dstUUID,
		DstTimestamp: *dstTimestamp,
	}

	return nil
}
