package authorityservice

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/testutil"
	"github.com/horizen-pes/pkg/storage/mockdb"
	"github.com/stretchr/testify/require"
)

func newTestService(t *testing.T, chainID uint64, ttl time.Duration, fixedTime time.Time) *AuthorityService {
	t.Helper()
	dl := mockdb.NewMockDataLayer()
	svc, err := NewAuthorityService(dl, chainID, ttl)
	require.NoError(t, err)
	svc.secret = bytes.Repeat([]byte{0x01}, 32)
	svc.clock = func() time.Time { return fixedTime }
	return svc
}

func signRequest(t *testing.T, chainID uint64, appID common.ApplicationIdType, reportID common.RequestIdType, nonce []byte) (string, ethCommon.Address) {
	t.Helper()
	key, err := ethCrypto.GenerateKey()
	require.NoError(t, err)
	hash := ethCrypto.Keccak256Hash(buildMessage(chainID, appID, reportID, nonce))
	sig, err := ethCrypto.Sign(hash.Bytes(), key)
	require.NoError(t, err)
	return hex.EncodeToString(sig), ethCrypto.PubkeyToAddress(key.PublicKey)
}

func TestHandleGetReportSuccess(t *testing.T) {
	chainID := uint64(42)
	now := time.Unix(1_700_000_000, 0)
	svc := newTestService(t, chainID, time.Minute, now)
	dl := svc.dataLayer.(*mockdb.MockDataLayer)

	saltHex := "00112233445566778899aabbccddeeff"
	saltBytes, _ := hex.DecodeString(saltHex)
	ts := now.Unix()
	nonceBytes := svc.computeNonce(saltBytes, ts)

	reportID := testutil.GenerateRandomRequestID()
	appID := common.NewApplicationId(1)
	signatureHex, authority := signRequest(t, chainID, appID, reportID, nonceBytes)

	report := &common.DeanonymizationReport{
		ApplicationID:   appID,
		ReportID:        reportID,
		Authority:       authority,
		EncryptedReport: []byte("encrypted-report"),
	}
	require.NoError(t, dl.StoreDeanonymizationReport(context.Background(), report))

	body := getReportRequest{
		ChainID:   chainID,
		AppID:     uint64(appID),
		ReportID:  reportID.String(),
		Salt:      saltHex,
		Nonce:     hex.EncodeToString(nonceBytes),
		Timestamp: ts,
		Signature: signatureHex,
	}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/getreport", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	svc.handleGetReport(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp getReportResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, authority.Hex(), resp.Authority)
	require.Equal(t, reportID.String(), resp.ReportID)
	require.Equal(t, appID.String(), resp.ApplicationID)
	require.Equal(t, hex.EncodeToString(report.EncryptedReport), resp.EncryptedReport)
}

func TestHandleGetReportExpiredNonce(t *testing.T) {
	chainID := uint64(42)
	now := time.Unix(1_700_000_000, 0)
	svc := newTestService(t, chainID, time.Minute, now)
	dl := svc.dataLayer.(*mockdb.MockDataLayer)

	saltHex := "00112233445566778899aabbccddeeff"
	saltBytes, _ := hex.DecodeString(saltHex)
	expiredTS := now.Add(-2 * time.Minute).Unix()
	nonceBytes := svc.computeNonce(saltBytes, expiredTS)

	reportID := testutil.GenerateRandomRequestID()
	appID := common.NewApplicationId(1)
	signatureHex, authority := signRequest(t, chainID, appID, reportID, nonceBytes)

	report := &common.DeanonymizationReport{
		ApplicationID:   appID,
		ReportID:        reportID,
		Authority:       authority,
		EncryptedReport: []byte("encrypted-report"),
	}
	require.NoError(t, dl.StoreDeanonymizationReport(context.Background(), report))

	body := getReportRequest{
		ChainID:   chainID,
		AppID:     uint64(appID),
		ReportID:  reportID.String(),
		Salt:      saltHex,
		Nonce:     hex.EncodeToString(nonceBytes),
		Timestamp: expiredTS,
		Signature: signatureHex,
	}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/getreport", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	svc.handleGetReport(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Contains(t, rr.Body.String(), "nonce expired")
}

func TestHandleGetReportSignatureMismatch(t *testing.T) {
	chainID := uint64(42)
	now := time.Unix(1_700_000_000, 0)
	svc := newTestService(t, chainID, time.Minute, now)
	dl := svc.dataLayer.(*mockdb.MockDataLayer)

	saltHex := "00112233445566778899aabbccddeeff"
	saltBytes, _ := hex.DecodeString(saltHex)
	ts := now.Unix()
	nonceBytes := svc.computeNonce(saltBytes, ts)

	reportID := testutil.GenerateRandomRequestID()
	appID := common.NewApplicationId(1)
	_, authority := signRequest(t, chainID, appID, reportID, nonceBytes)

	// Tamper signature to make recovery fail
	badSignatureHex := hex.EncodeToString(bytes.Repeat([]byte{0}, 65))

	report := &common.DeanonymizationReport{
		ApplicationID:   appID,
		ReportID:        reportID,
		Authority:       authority,
		EncryptedReport: []byte("encrypted-report"),
	}
	require.NoError(t, dl.StoreDeanonymizationReport(context.Background(), report))

	body := getReportRequest{
		ChainID:   chainID,
		AppID:     uint64(appID),
		ReportID:  reportID.String(),
		Salt:      saltHex,
		Nonce:     hex.EncodeToString(nonceBytes),
		Timestamp: ts,
		Signature: badSignatureHex,
	}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/getreport", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	svc.handleGetReport(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Contains(t, rr.Body.String(), "failed to recover signer")
}

func TestHandleGetReportAuthorityMismatch(t *testing.T) {
	chainID := uint64(42)
	now := time.Unix(1_700_000_000, 0)
	svc := newTestService(t, chainID, time.Minute, now)
	dl := svc.dataLayer.(*mockdb.MockDataLayer)

	saltHex := "00112233445566778899aabbccddeeff"
	saltBytes, _ := hex.DecodeString(saltHex)
	ts := now.Unix()
	nonceBytes := svc.computeNonce(saltBytes, ts)

	reportID := testutil.GenerateRandomRequestID()
	appID := common.NewApplicationId(1)
	signatureHex, _ := signRequest(t, chainID, appID, reportID, nonceBytes)

	report := &common.DeanonymizationReport{
		ApplicationID:   appID,
		ReportID:        reportID,
		Authority:       ethCommon.HexToAddress("0x0000000000000000000000000000000000000001"),
		EncryptedReport: []byte("encrypted-report"),
	}
	require.NoError(t, dl.StoreDeanonymizationReport(context.Background(), report))

	body := getReportRequest{
		ChainID:   chainID,
		AppID:     uint64(appID),
		ReportID:  reportID.String(),
		Salt:      saltHex,
		Nonce:     hex.EncodeToString(nonceBytes),
		Timestamp: ts,
		Signature: signatureHex,
	}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/getreport", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	svc.handleGetReport(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
	require.Contains(t, rr.Body.String(), "authority mismatch")
}

func TestHandleGetReportAppMismatch(t *testing.T) {
	chainID := uint64(42)
	now := time.Unix(1_700_000_000, 0)
	svc := newTestService(t, chainID, time.Minute, now)
	dl := svc.dataLayer.(*mockdb.MockDataLayer)

	saltHex := "00112233445566778899aabbccddeeff"
	saltBytes, _ := hex.DecodeString(saltHex)
	ts := now.Unix()
	nonceBytes := svc.computeNonce(saltBytes, ts)

	reportID := testutil.GenerateRandomRequestID()
	appID := common.NewApplicationId(2)
	signatureHex, authority := signRequest(t, chainID, appID, reportID, nonceBytes)

	report := &common.DeanonymizationReport{
		ApplicationID:   common.NewApplicationId(3), // mismatch
		ReportID:        reportID,
		Authority:       authority,
		EncryptedReport: []byte("encrypted-report"),
	}
	require.NoError(t, dl.StoreDeanonymizationReport(context.Background(), report))

	body := getReportRequest{
		ChainID:   chainID,
		AppID:     uint64(appID),
		ReportID:  reportID.String(),
		Salt:      saltHex,
		Nonce:     hex.EncodeToString(nonceBytes),
		Timestamp: ts,
		Signature: signatureHex,
	}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/getreport", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	svc.handleGetReport(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
	require.Contains(t, rr.Body.String(), "applicationId mismatch")
}
