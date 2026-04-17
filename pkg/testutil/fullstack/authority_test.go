package fullstack_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/HorizenOfficial/vela/pkg/authorityservice/api"
	"github.com/HorizenOfficial/vela/pkg/common"
	commontestutil "github.com/HorizenOfficial/vela/pkg/common/testutil"
	"github.com/HorizenOfficial/vela/pkg/executor"
	"github.com/HorizenOfficial/vela/pkg/logger"
	"github.com/HorizenOfficial/vela/pkg/manager"
	"github.com/HorizenOfficial/vela/pkg/testutil/fullstack"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

// TestAuthorityServiceNonceEndpoint is a minimal boot check: it proves the
// in-process authority server binds, serves GET /nonce, and returns a
// well-formed NonceResponse. No subgraph or report state is required.
func TestAuthorityServiceNonceEndpoint(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("skipping fullstack test under CI_FLAG")
	}

	suite := newAuthoritySuite(t)
	defer func() { _ = suite.Cleanup() }()

	resp, err := http.Get(suite.GetAuthorityServiceURL() + "/nonce")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var nonceResp api.NonceResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&nonceResp))
	require.NotEmpty(t, nonceResp.Salt, "salt should be populated")
	require.NotEmpty(t, nonceResp.Nonce, "nonce should be populated")
	require.NotZero(t, nonceResp.Timestamp, "timestamp should be populated")

	_, err = hex.DecodeString(nonceResp.Salt)
	require.NoError(t, err, "salt should be hex")
	_, err = hex.DecodeString(nonceResp.Nonce)
	require.NoError(t, err, "nonce should be hex")

	require.NotEqual(t, ethCommon.Address{}, suite.GetAuthorityAddress())
}

// TestAuthorityServiceGetReportRoundTrip exercises the full /nonce + /getreport
// path end-to-end against an in-process authority, subgraph, and filesystem
// reports dir. It synthesizes the report file and subgraph completion directly
// (via InjectRequestCompleted) since Phase 3 does not drive a real deanonymize
// request through the manager.
func TestAuthorityServiceGetReportRoundTrip(t *testing.T) {
	if os.Getenv("CI_FLAG") != "" {
		t.Skip("skipping fullstack test under CI_FLAG")
	}

	suite := newAuthoritySuite(t)
	defer func() { _ = suite.Cleanup() }()

	appID := common.NewApplicationId(1)
	reportID := commontestutil.GenerateRandomRequestID()
	encryptedBody := []byte("encrypted-report-body")

	writeReport(t, suite.GetReportsPath(), &common.DeanonymizationReport{
		ApplicationID:   appID,
		ReportID:        reportID,
		Authority:       suite.GetAuthorityAddress(),
		EncryptedReport: encryptedBody,
	})

	// Inject a synthetic on-chain completion so the /getreport handler's
	// subgraph check passes.
	suite.GetSubgraph().InjectRequestCompleted(appID, reportID, common.RequestResultOK)

	// Step 1: GET /nonce
	nonceResp := fetchNonce(t, suite.GetAuthorityServiceURL())

	nonceBytes, err := hex.DecodeString(nonceResp.Nonce)
	require.NoError(t, err)

	// Step 2: sign the canonical message with the authority's key
	msg := api.BuildMessage(fullstack.ChainID.Uint64(), appID, reportID, nonceBytes)
	digest := ethCrypto.Keccak256Hash(msg)
	sig, err := ethCrypto.Sign(digest.Bytes(), suite.GetAuthorityKey())
	require.NoError(t, err)

	// Step 3: POST /getreport
	reqBody := api.GetReportRequest{
		ChainID:   fullstack.ChainID.Uint64(),
		AppID:     uint64(appID),
		ReportID:  reportID.String(),
		Salt:      nonceResp.Salt,
		Nonce:     nonceResp.Nonce,
		Timestamp: nonceResp.Timestamp,
		Signature: hex.EncodeToString(sig),
	}
	payload, err := json.Marshal(reqBody)
	require.NoError(t, err)

	resp, err := http.Post(suite.GetAuthorityServiceURL()+"/getreport", "application/json", bytes.NewReader(payload))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "expected 200 from /getreport")

	var got api.GetReportResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	require.Equal(t, appID.String(), got.ApplicationID)
	require.Equal(t, reportID.String(), got.ReportID)
	require.Equal(t, suite.GetAuthorityAddress().Hex(), got.Authority)
	require.Equal(t, hex.EncodeToString(encryptedBody), got.EncryptedReport)
}

// --- helpers ---

func newAuthoritySuite(t *testing.T) *fullstack.FullStackSystemTestSuite {
	t.Helper()
	// manager.LoadConfig requires MANAGER_ARTIFACTS_PATH. TestSuiteCore overrides
	// this with a temp dir during construction — we just need LoadConfig to pass.
	t.Setenv("MANAGER_ARTIFACTS_PATH", t.TempDir())
	// Select unsafe keyset recovery (type 0) so GenerateEnclaveKeySet does not
	// need a KMS client. Matches what system_tests/executor.conf does.
	t.Setenv("EXECUTOR_KEYSET_RECOVERY_TYPE", "0")
	mgrCfg, err := manager.LoadConfig()
	require.NoError(t, err)
	execCfg, err := executor.LoadConfig()
	require.NoError(t, err)
	keySet, recovery, err := executor.GenerateEnclaveKeySet(t.Context(), execCfg.KeySetRecoveryType, nil, nil, "")
	require.NoError(t, err)
	logCfg := &logger.Config{Kind: "zerolog", Console: true, ConsoleLevel: "info"}
	return fullstack.NewFullStackSystemTestSuiteWithConfigs(t, "mock-runtime", mgrCfg, execCfg, keySet, recovery, logCfg, logCfg)
}

func fetchNonce(t *testing.T, baseURL string) api.NonceResponse {
	t.Helper()
	resp, err := http.Get(baseURL + "/nonce")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var nr api.NonceResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&nr))
	return nr
}

func writeReport(t *testing.T, reportsDir string, report *common.DeanonymizationReport) {
	t.Helper()
	require.NoError(t, os.MkdirAll(reportsDir, 0o755))
	path := filepath.Join(reportsDir, common.ReportFilename(report.ApplicationID, report.ReportID))
	data, err := json.Marshal(report)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, data, 0o644))
}
