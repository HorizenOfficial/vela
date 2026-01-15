package subgraph

import (
	"context"
	"strings"
	"testing"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/testutil"
	"github.com/horizen-pes/pkg/crypto"
	"github.com/stretchr/testify/require"
)

func TestFetchAndDecryptUserEvents_StopAtFirst(t *testing.T) {
	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	appID := common.NewApplicationId(1)
	reqID1 := testutil.GenerateRandomRequestID()
	reqID2 := testutil.GenerateRandomRequestID()

	ev1Cipher, err := crypto.Encrypt(teeKey, userKey.PublicKey(), []byte("msg-1"))
	require.NoError(t, err)
	ev2Cipher, err := crypto.Encrypt(teeKey, userKey.PublicKey(), []byte("msg-2"))
	require.NoError(t, err)

	mock := NewMockClient().WithUserEvents(appID, []UserEvent{
		{ApplicationID: appID, RequestID: reqID2, EncryptedData: ev2Cipher, EventSubType: "a", BlockNumber: 3},
		{ApplicationID: appID, RequestID: reqID1, EncryptedData: ev1Cipher, EventSubType: "a", BlockNumber: 2},
	})

	result, err := FetchAndDecryptUserEvents(context.Background(), mock, teeKey.PublicKey(), *userKey, appID, "", 10, nil, true)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, []byte("msg-2"), result[0])
}

func TestFetchAndDecryptUserEvents_Filter(t *testing.T) {
	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	appID := common.NewApplicationId(2)
	reqID1 := testutil.GenerateRandomRequestID()
	reqID2 := testutil.GenerateRandomRequestID()

	ev1Cipher, err := crypto.Encrypt(teeKey, userKey.PublicKey(), []byte("keep-this"))
	require.NoError(t, err)
	ev2Cipher, err := crypto.Encrypt(teeKey, userKey.PublicKey(), []byte("drop-this"))
	require.NoError(t, err)

	mock := NewMockClient().WithUserEvents(appID, []UserEvent{
		{ApplicationID: appID, RequestID: reqID1, EncryptedData: ev1Cipher, EventSubType: "b", BlockNumber: 5},
		{ApplicationID: appID, RequestID: reqID2, EncryptedData: ev2Cipher, EventSubType: "b", BlockNumber: 4},
	})

	filter := func(data []byte) bool {
		return strings.Contains(string(data), "keep")
	}

	result, err := FetchAndDecryptUserEvents(context.Background(), mock, teeKey.PublicKey(), *userKey, appID, "", 10, filter, false)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, []byte("keep-this"), result[0])
}
