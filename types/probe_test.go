package types_test

import (
	"testing"
	"time"

	"github.com/flashbots/latency-monitor/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestProbeEncodeDecode(t *testing.T) {
	pOrg := types.Probe{
		Sequence:     42,
		SrcUUID:      uuid.New(),
		SrcTimestamp: time.Now(),
		DstUUID:      uuid.New(),
		DstTimestamp: time.Now(),
	}

	b, err := pOrg.MarshalBinary()
	require.NoError(t, err)

	pRes := &types.Probe{}
	err = pRes.UnmarshalBinary(b)
	require.NoError(t, err)

	require.Equal(t, pOrg.Sequence, pRes.Sequence)
	require.Equal(t, pOrg.SrcUUID, pRes.SrcUUID)
	require.Equal(t, pOrg.SrcTimestamp.UnixNano(), pRes.SrcTimestamp.UnixNano()) // otherwise, monotonic clock will drift
	require.Equal(t, pOrg.DstUUID, pRes.DstUUID)
	require.Equal(t, pOrg.DstTimestamp.UnixNano(), pRes.DstTimestamp.UnixNano()) // otherwise, monotonic clock will drift
}
