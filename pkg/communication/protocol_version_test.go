package communication

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsCompatible(t *testing.T) {
	require.True(t, IsCompatible(WireProtocolVersion), "same version must be compatible")
	require.False(t, IsCompatible(WireProtocolVersion+1), "a higher version is incompatible")
	require.False(t, IsCompatible(0), "version 0 (legacy/absent) is incompatible with a non-zero WireProtocolVersion")
}

// TestProtocolVersion_AbsentFieldIsZero documents the bootstrap behavior: a peer
// that predates the field sends JSON without it, which unmarshals to version 0.
func TestProtocolVersion_AbsentFieldIsZero(t *testing.T) {
	var reqData GetKeysetRecoveryRequestData
	require.NoError(t, json.Unmarshal([]byte(`{}`), &reqData))
	require.Equal(t, uint32(0), reqData.WireProtocolVersion)

	var respData GetKeysetRecoveryResponseData
	require.NoError(t, json.Unmarshal([]byte(`{"dataFound":false}`), &respData))
	require.Equal(t, uint32(0), respData.WireProtocolVersion)
}

func TestIncompatibleProtocolError_Message(t *testing.T) {
	err := &IncompatibleProtocolError{Local: 1, Peer: 0}
	require.Contains(t, err.Error(), "local=1")
	require.Contains(t, err.Error(), "peer=0")
}
