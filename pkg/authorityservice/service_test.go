package authorityservice

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/horizen-pes/pkg/authorityservice/api"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/testutil"
	"github.com/horizen-pes/pkg/logger"
	"github.com/stretchr/testify/require"
)

func newTestService(t *testing.T, chainID uint64, ttl time.Duration, fixedTime time.Time) *AuthorityService {
	t.Helper()
	dir := t.TempDir()
	testLogger := logger.NewLogger(
		&logger.Config{
			Kind:         "zerolog",
			ConsoleColor: false, // colors can print escape chars on tty
			Console:      true,
			ConsoleLevel: "trace",
			//FileName:     "qqq.log",
			//FileLevel:    "info",
		},
	)
	svc, err := NewAuthorityService(chainID, ttl, dir, testLogger)
	require.NoError(t, err)
	svc.secret = bytes.Repeat([]byte{0x01}, 32)
	svc.clock = func() time.Time { return fixedTime }
	return svc
}

func signRequest(t *testing.T, chainID uint64, appID common.ApplicationIdType, reportID common.RequestIdType, nonce []byte) (string, ethCommon.Address) {
	t.Helper()
	key, err := ethCrypto.GenerateKey()
	require.NoError(t, err)
	hash := ethCrypto.Keccak256Hash(api.BuildMessage(chainID, appID, reportID, nonce))
	sig, err := ethCrypto.Sign(hash.Bytes(), key)
	require.NoError(t, err)
	return hex.EncodeToString(sig), ethCrypto.PubkeyToAddress(key.PublicKey)
}

func writeReport(t *testing.T, svc *AuthorityService, report *common.DeanonymizationReport) {
	t.Helper()
	if err := os.MkdirAll(svc.reportPath, 0755); err != nil {
		t.Fatalf("failed to create report dir: %v", err)
	}
	filePath := filepath.Join(svc.reportPath, common.ReportFilename(report.ApplicationID, report.ReportID))
	data, err := json.Marshal(report)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filePath, data, 0644))
}

func TestHandleGetReportSuccess(t *testing.T) {
	chainID := uint64(42)
	now := time.Unix(1_700_000_000, 0)
	svc := newTestService(t, chainID, time.Minute, now)

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
	writeReport(t, svc, report)

	body := api.GetReportRequest{
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

	var resp api.GetReportResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, authority.Hex(), resp.Authority)
	require.Equal(t, reportID.String(), resp.ReportID)
	require.Equal(t, appID.String(), resp.ApplicationID)
	require.Equal(t, hex.EncodeToString(report.EncryptedReport), resp.EncryptedReport)
}

func TestHandleGetReportUnexpectedChainID(t *testing.T) {
	chainID := uint64(42)
	now := time.Unix(1_700_000_000, 0)
	svc := newTestService(t, chainID, time.Minute, now)

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
	writeReport(t, svc, report)

	body := api.GetReportRequest{
		ChainID:   chainID + 1, // wrong chain
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

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "bad request")
}

func TestHandleGetReportReportNotFound(t *testing.T) {
	chainID := uint64(42)
	now := time.Unix(1_700_000_000, 0)
	svc := newTestService(t, chainID, time.Minute, now)

	saltHex := "00112233445566778899aabbccddeeff"
	saltBytes, _ := hex.DecodeString(saltHex)
	ts := now.Unix()
	nonceBytes := svc.computeNonce(saltBytes, ts)

	reportID := testutil.GenerateRandomRequestID()
	appID := common.NewApplicationId(1)
	signatureHex, _ := signRequest(t, chainID, appID, reportID, nonceBytes)

	body := api.GetReportRequest{
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

	require.Equal(t, http.StatusNotFound, rr.Code)
	require.Contains(t, rr.Body.String(), "not found")
}

func TestHandleGetReportFutureNonce(t *testing.T) {
	chainID := uint64(42)
	now := time.Unix(1_700_000_000, 0)
	svc := newTestService(t, chainID, time.Minute, now)

	saltHex := "00112233445566778899aabbccddeeff"
	saltBytes, _ := hex.DecodeString(saltHex)
	ts := now.Add(2 * time.Minute).Unix()
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
	writeReport(t, svc, report)

	body := api.GetReportRequest{
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

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Contains(t, rr.Body.String(), "unauthorized")
}

func TestHandleGetReportInvalidNonceContent(t *testing.T) {
	chainID := uint64(42)
	now := time.Unix(1_700_000_000, 0)
	svc := newTestService(t, chainID, time.Minute, now)

	saltHex := "00112233445566778899aabbccddeeff"
	saltBytes, _ := hex.DecodeString(saltHex)
	ts := now.Unix()
	wrongTime := time.Unix(1_760_000_000, 0).Unix()
	// compute a valid nonce then tamper it
	badNonce := svc.computeNonce(saltBytes, wrongTime)

	reportID := testutil.GenerateRandomRequestID()
	appID := common.NewApplicationId(1)
	signatureHex, authority := signRequest(t, chainID, appID, reportID, badNonce)

	report := &common.DeanonymizationReport{
		ApplicationID:   appID,
		ReportID:        reportID,
		Authority:       authority,
		EncryptedReport: []byte("encrypted-report"),
	}
	writeReport(t, svc, report)

	body := api.GetReportRequest{
		ChainID:   chainID,
		AppID:     uint64(appID),
		ReportID:  reportID.String(),
		Salt:      saltHex,
		Nonce:     hex.EncodeToString(badNonce),
		Timestamp: ts,
		Signature: signatureHex,
	}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/getreport", bytes.NewReader(payload))
	rr := httptest.NewRecorder()

	svc.handleGetReport(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
	require.Contains(t, rr.Body.String(), "unauthorized")
}

func TestHandleGetReportBadMethod(t *testing.T) {
	svc := newTestService(t, 0, time.Minute, time.Unix(1_700_000_000, 0))

	req := httptest.NewRequest(http.MethodGet, "/getreport", nil)
	rr := httptest.NewRecorder()

	svc.handleGetReport(rr, req)

	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	require.Contains(t, rr.Body.String(), "method not allowed")
}

func TestHandleNonceBadMethod(t *testing.T) {
	svc := newTestService(t, 0, time.Minute, time.Unix(1_700_000_000, 0))

	req := httptest.NewRequest(http.MethodPost, "/nonce", nil)
	rr := httptest.NewRecorder()

	svc.handleNonce(rr, req)

	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	require.Contains(t, rr.Body.String(), "method not allowed")
}

func TestHandleGetReportExpiredNonce(t *testing.T) {
	chainID := uint64(42)
	now := time.Unix(1_700_000_000, 0)
	svc := newTestService(t, chainID, time.Minute, now)

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
	writeReport(t, svc, report)

	body := api.GetReportRequest{
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
	require.Contains(t, rr.Body.String(), "unauthorized")
}

func TestHandleGetReportSignatureMismatch(t *testing.T) {
	chainID := uint64(42)
	now := time.Unix(1_700_000_000, 0)
	svc := newTestService(t, chainID, time.Minute, now)

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
	writeReport(t, svc, report)

	body := api.GetReportRequest{
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
	require.Contains(t, rr.Body.String(), "unauthorized")
}

func TestHandleGetReportAuthorityMismatch(t *testing.T) {
	chainID := uint64(42)
	now := time.Unix(1_700_000_000, 0)
	svc := newTestService(t, chainID, time.Minute, now)

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
	writeReport(t, svc, report)

	body := api.GetReportRequest{
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
	require.Contains(t, rr.Body.String(), "forbidden")
}

func TestHandleGetReportAppMismatch(t *testing.T) {
	chainID := uint64(42)
	now := time.Unix(1_700_000_000, 0)
	svc := newTestService(t, chainID, time.Minute, now)

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
	writeReport(t, svc, report)

	body := api.GetReportRequest{
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
	require.Contains(t, rr.Body.String(), "forbidden")
}
