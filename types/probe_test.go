package types_test

import (
	"testing"
	"time"

	"github.com/flashbots/latency-monitor/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestProbeEncodeDecode(t *testing.T) {
	srcLocation := types.Location{}
	dstLocation := types.Location{}

	copy(srcLocation[:], []byte("sourceLocation"))
	copy(dstLocation[:], []byte("destinationLocation"))

	pOrg := types.Probe{
		Sequence:     42,
		SrcUUID:      uuid.New(),
		SrcTimestamp: time.Now(),
		SrcLocation:  types.Location(srcLocation),
		DstUUID:      uuid.New(),
		DstTimestamp: time.Now(),
		DstLocation:  types.Location(dstLocation),
	}

	b, err := pOrg.MarshalBinary()
	require.NoError(t, err)

	pRes := &types.Probe{}
	err = pRes.UnmarshalBinary(b)
	require.NoError(t, err)

	require.Equal(t, pOrg.Sequence, pRes.Sequence)
	require.Equal(t, pOrg.SrcUUID, pRes.SrcUUID)
	require.Equal(t, pOrg.SrcTimestamp.UnixNano(), pRes.SrcTimestamp.UnixNano()) // otherwise, monotonic clock will drift
	require.Equal(t, pOrg.SrcLocation, pRes.SrcLocation)
	require.Equal(t, pOrg.DstUUID, pRes.DstUUID)
	require.Equal(t, pOrg.DstTimestamp.UnixNano(), pRes.DstTimestamp.UnixNano()) // otherwise, monotonic clock will drift
	require.Equal(t, pOrg.DstLocation, pRes.DstLocation)

	t.Logf("Src: %s", pRes.SrcLocation.String())
	t.Logf("Dst: %s", pRes.DstLocation.String())
}
