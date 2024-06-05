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
	SrcLocation  Location
	DstUUID      uuid.UUID
	DstTimestamp time.Time
	DstLocation  Location
}

func ProbeSize() int {
	return 142
}

var (
	ErrProbeFailedToEncodeBinaryRepresentation = errors.New("failed to encode probe into its binary representation")
	ErrProbeFailedToDecodeBinaryRepresentation = errors.New("failed to decode probe from its binary representation")
)

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

	data := make([]byte, ProbeSize())

	binary.LittleEndian.PutUint64(data[0:8], p.Sequence) // 000..007  : 8 bytes
	copy(data[8:24], p.SrcUUID[:])                       // 008..023  : 16 bytes
	copy(data[24:39], rawSrcTimestamp)                   // 024..038  : 15 bytes
	copy(data[39:75], p.SrcLocation[:])                  // 039..074  : 36 bytes
	copy(data[75:91], p.DstUUID[:])                      // 075..090  : 16 bytes
	copy(data[91:106], rawDstTimestamp)                  // 091..105  : 15 bytes
	copy(data[106:142], p.DstLocation[:])                // 106..142  : 36 bytes

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

	srcLocation := Location{}
	copy(srcLocation[:], data[39:75])

	dstUUID, err := uuid.FromBytes(data[75:91])
	if err != nil {
		return fmt.Errorf("%w: DstUUID: %w",
			ErrProbeFailedToDecodeBinaryRepresentation, err,
		)
	}

	dstTimestamp := &time.Time{}
	if err := dstTimestamp.UnmarshalBinary(data[91:106]); err != nil {
		return fmt.Errorf("%w: DstTimestamp: %w",
			ErrProbeFailedToDecodeBinaryRepresentation, err,
		)
	}

	dstLocation := Location{}
	copy(dstLocation[:], data[106:142])

	*p = Probe{
		Sequence:     binary.LittleEndian.Uint64(data[:8]),
		SrcUUID:      srcUUID,
		SrcTimestamp: *srcTimestamp,
		SrcLocation:  srcLocation,
		DstUUID:      dstUUID,
		DstTimestamp: *dstTimestamp,
		DstLocation:  dstLocation,
	}

	return nil
}
