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

func withUserEventsPageSize(t *testing.T, size int) {
	t.Helper()
	old := userEventsPageSize
	userEventsPageSize = size
	t.Cleanup(func() {
		userEventsPageSize = old
	})
}

func TestFetchAndDecryptUserEvents_LimitOne(t *testing.T) {
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

	result, err := FetchAndDecryptUserEvents(context.Background(), mock, teeKey.PublicKey(), *userKey, appID, "", 1, nil)
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

	result, err := FetchAndDecryptUserEvents(context.Background(), mock, teeKey.PublicKey(), *userKey, appID, "", 10, filter)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, []byte("keep-this"), result[0])
}

func TestFetchAndDecryptUserEvents_UserSpecificDecryption(t *testing.T) {
	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	userKeyA, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	userKeyB, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	appID := common.NewApplicationId(3)
	reqID1 := testutil.GenerateRandomRequestID()
	reqID2 := testutil.GenerateRandomRequestID()

	ev1Cipher, err := crypto.Encrypt(teeKey, userKeyA.PublicKey(), []byte("user-a"))
	require.NoError(t, err)
	ev2Cipher, err := crypto.Encrypt(teeKey, userKeyB.PublicKey(), []byte("user-b"))
	require.NoError(t, err)

	mock := NewMockClient().WithUserEvents(appID, []UserEvent{
		{ApplicationID: appID, RequestID: reqID1, EncryptedData: ev1Cipher, EventSubType: "c", BlockNumber: 3},
		{ApplicationID: appID, RequestID: reqID2, EncryptedData: ev2Cipher, EventSubType: "c", BlockNumber: 2},
	})

	result, err := FetchAndDecryptUserEvents(context.Background(), mock, teeKey.PublicKey(), *userKeyA, appID, "", 10, nil)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, []byte("user-a"), result[0])

	result, err = FetchAndDecryptUserEvents(context.Background(), mock, teeKey.PublicKey(), *userKeyB, appID, "", 10, nil)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, []byte("user-b"), result[0])
}

func TestFetchAndDecryptUserEvents_PaginatesUntilMatch(t *testing.T) {
	withUserEventsPageSize(t, 1)

	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	appID := common.NewApplicationId(3)
	reqID1 := testutil.GenerateRandomRequestID()
	reqID2 := testutil.GenerateRandomRequestID()

	ev1Cipher, err := crypto.Encrypt(teeKey, userKey.PublicKey(), []byte("skip-me"))
	require.NoError(t, err)
	ev2Cipher, err := crypto.Encrypt(teeKey, userKey.PublicKey(), []byte("target"))
	require.NoError(t, err)

	mock := NewMockClient().WithUserEvents(appID, []UserEvent{
		{ApplicationID: appID, RequestID: reqID1, EncryptedData: ev1Cipher, EventSubType: "c", BlockNumber: 2},
		{ApplicationID: appID, RequestID: reqID2, EncryptedData: ev2Cipher, EventSubType: "c", BlockNumber: 1},
	})

	filter := func(data []byte) bool {
		return strings.Contains(string(data), "target")
	}

	result, err := FetchAndDecryptUserEvents(context.Background(), mock, teeKey.PublicKey(), *userKey, appID, "", 1, filter)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, []byte("target"), result[0])
}

func TestFetchAndDecryptUserEvents_MaxResults(t *testing.T) {
	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	appID := common.NewApplicationId(4)
	reqID1 := testutil.GenerateRandomRequestID()
	reqID2 := testutil.GenerateRandomRequestID()

	ev1Cipher, err := crypto.Encrypt(teeKey, userKey.PublicKey(), []byte("first"))
	require.NoError(t, err)
	ev2Cipher, err := crypto.Encrypt(teeKey, userKey.PublicKey(), []byte("second"))
	require.NoError(t, err)

	mock := NewMockClient().WithUserEvents(appID, []UserEvent{
		{ApplicationID: appID, RequestID: reqID1, EncryptedData: ev1Cipher, EventSubType: "d", BlockNumber: 2},
		{ApplicationID: appID, RequestID: reqID2, EncryptedData: ev2Cipher, EventSubType: "d", BlockNumber: 1},
	})

	result, err := FetchAndDecryptUserEvents(context.Background(), mock, teeKey.PublicKey(), *userKey, appID, "", 1, nil)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, []byte("first"), result[0])
}

func TestFetchAndDecryptUserEvents_NoLimit(t *testing.T) {
	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	appID := common.NewApplicationId(5)
	reqID1 := testutil.GenerateRandomRequestID()
	reqID2 := testutil.GenerateRandomRequestID()

	ev1Cipher, err := crypto.Encrypt(teeKey, userKey.PublicKey(), []byte("one"))
	require.NoError(t, err)
	ev2Cipher, err := crypto.Encrypt(teeKey, userKey.PublicKey(), []byte("two"))
	require.NoError(t, err)

	mock := NewMockClient().WithUserEvents(appID, []UserEvent{
		{ApplicationID: appID, RequestID: reqID1, EncryptedData: ev1Cipher, EventSubType: "e", BlockNumber: 2},
		{ApplicationID: appID, RequestID: reqID2, EncryptedData: ev2Cipher, EventSubType: "e", BlockNumber: 1},
	})

	result, err := FetchAndDecryptUserEvents(context.Background(), mock, teeKey.PublicKey(), *userKey, appID, "", 0, nil)
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, []byte("one"), result[0])
	require.Equal(t, []byte("two"), result[1])
}

func TestFetchAndDecryptUserEvents_OrderWithinBlock(t *testing.T) {
	teeKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)
	userKey, err := crypto.GeneratePrivateKeyP521()
	require.NoError(t, err)

	appID := common.NewApplicationId(6)
	reqID1 := testutil.GenerateRandomRequestID()
	reqID2 := testutil.GenerateRandomRequestID()

	ev1Cipher, err := crypto.Encrypt(teeKey, userKey.PublicKey(), []byte("first-log"))
	require.NoError(t, err)
	ev2Cipher, err := crypto.Encrypt(teeKey, userKey.PublicKey(), []byte("second-log"))
	require.NoError(t, err)

	mock := NewMockClient().WithUserEvents(appID, []UserEvent{
		{ApplicationID: appID, RequestID: reqID1, EncryptedData: ev1Cipher, EventSubType: "f", BlockNumber: 10, LogIndex: 1},
		{ApplicationID: appID, RequestID: reqID2, EncryptedData: ev2Cipher, EventSubType: "f", BlockNumber: 10, LogIndex: 2},
	})

	result, err := FetchAndDecryptUserEvents(context.Background(), mock, teeKey.PublicKey(), *userKey, appID, "", 1, nil)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, []byte("second-log"), result[0])
}
